// Package programs provides the programs service for the host application.
package programs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type ProgramService struct {
	db         *ProgramDB
	rawgClient *RAWGClient
}

// NewService creates a new program service
func NewService(
	dbPath string,
	rawgAPIKey string,
) *ProgramService {
	db, err := NewProgramDB(dbPath)
	if err != nil {
		//TODO: maybe we should panic
		log.Printf("Failed to initialize program database: %v", err)
	}
	return &ProgramService{
		db:         db,
		rawgClient: NewRAWGClient(rawgAPIKey),
	}
}

// DiscoverAndSavePrograms discovers program in common and custom paths and saves them to the local db
func (s *ProgramService) DiscoverAndSavePrograms(paths []string) error {
	programs, err := s.DiscoverPrograms()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	for _, program := range programs {
		if _, ok := seen[program.Path]; ok {
			continue
		}

		program = s.RawgSearch(program)

		if s.db != nil {
			if err := s.db.SaveProgram(&program); err != nil {
				log.Printf("Failed to save program %s: %v", program.Name, err)
			}
		}

		seen[program.Path] = true
	}

	if len(paths) > 0 {
		customPathProgs, err := s.DiscoverProgramsIn(paths)
		if err != nil {
			return err
		}
		for _, program := range customPathProgs {
			if _, ok := seen[program.Path]; ok {
				continue
			}

			program = s.RawgSearch(program)

			if s.db != nil {
				if err := s.db.SaveProgram(&program); err != nil {
					log.Printf("Failed to save program %s: %v", program.Name, err)
				}
			}

			seen[program.Path] = true
		}
	}

	if s.db != nil {
		if err := s.db.CleanupNonExistentPrograms(); err != nil {
			log.Printf("Failed to cleanup non-existent programs: %v", err)
		}
	}

	return nil
}

// GetLocalPrograms gets all programs from the local db
func (s *ProgramService) GetLocalPrograms() ([]*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}

	return s.db.GetPrograms()
}

// GetLocalProgramByPath gets a program from the local db by path
func (s *ProgramService) GetLocalProgramByPath(path string) (*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}
	return s.db.GetProgramByPath(path)
}

// DiscoverPrograms discovers programs in the common paths
func (s *ProgramService) DiscoverPrograms() ([]Program, error) {
	var programs []Program

	switch runtime.GOOS {
	case "windows":
		// Common game directories on Windows
		commonPaths := []string{
			"C:\\Program Files (x86)\\Steam\\steamapps\\common",
			"C:\\Program Files\\Steam\\steamapps\\common",
			"C:\\Program Files (x86)\\Epic Games",
			"C:\\Program Files\\Epic Games",
			"C:\\Games",
		}

		for _, basePath := range commonPaths {
			discoveredPrograms, err := s.scanDirectoryForPrograms(basePath)
			if err != nil {
				// Log error but continue with other paths
				continue
			}
			programs = append(programs, discoveredPrograms...)
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return programs, nil
}

// DiscoverProgramsIn discovers programs in the given paths
func (s *ProgramService) DiscoverProgramsIn(paths []string) ([]Program, error) {
	var programs []Program
	for _, p := range paths {
		discoveredPrograms, err := s.scanDirectoryForPrograms(p)
		if err != nil {
			continue
		}
		programs = append(programs, discoveredPrograms...)
	}
	return programs, nil
}

func (s *ProgramService) scanDirectoryForPrograms(path string) ([]Program, error) {
	var programs []Program

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if shgouldIgnoreProgram(info.Name()) {
			return nil
		}

		if info.Size() < 5*1024*1024 {
			return nil
		}

		switch runtime.GOOS {
		case "windows":
			if strings.HasSuffix(strings.ToLower(info.Name()), ".exe") {
				programs = append(programs, Program{
					Name: info.Name(),
					Path: p,
				})
			}
		default:
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return programs, nil
}

func (s *ProgramService) SaveProgram(req CreateProgramRequest) (*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}

	progr := &Program{
		Name:        req.Name,
		Path:        req.Path,
		Description: req.Description,
		HostID:      req.HostID,
	}

	if err := s.db.SaveProgram(progr); err != nil {
		return nil, err
	}

	return progr, nil
}

func (s *ProgramService) LaunchProgram(path string) (*exec.Cmd, error) {
	cmd := exec.Command(path)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to launch program: %w", err)
	}
	return cmd, nil
}

func (s *ProgramService) GetWindowTitle(processName string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Simplest way to do this????
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Get-Process | Where-Object {$_.ProcessName -like '*%s*'} | ForEach-Object { (Get-WindowTitle -ProcessId $_.Id) }", processName))
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func shgouldIgnoreProgram(name string) bool {
	var ignoredExecutables = []string{
		"unins", "uninstall", "setup", "update", "updater",
		"launcher", "crashreporter", "helper",
		"steamerrorreporter", "steamservice",
		"witcherscriptmerger", "quickbms", "scc", "kdiff3",
		"vcredist", "dotnet", "dxsetup", "crashreport", "reportclient",
		"helper", "bootstrapper", "redist", "redistributable",
	}

	base := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))

	for _, ignore := range ignoredExecutables {
		if strings.Contains(base, ignore) {
			return true
		}
	}
	return false
}

func cleanQueryFromExeName(name string) (string, bool) {
	base := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))

	// hard ignore
	if shgouldIgnoreProgram(base) {
		return "", false
	}

	// too short or mostly symbols
	if len(base) < 4 {
		return "", false
	}

	// strip known suffixes
	suffixes := []string{"launcher", "setup", "updater", "merger", "tool"}
	for _, s := range suffixes {
		if strings.Contains(base, s) {
			return "", false
		}
	}

	return base, true
}

func (s *ProgramService) RawgSearch(program Program) Program {
	cleanedName, ok := cleanQueryFromExeName(program.Name)
	if !ok {
		return program
	}
	results, err := s.rawgClient.SearchGame(cleanedName)
	log.Printf("Searching for program %s: %v", program.Name, results)

	if err == nil && len(results) > 0 {
		for _, game := range results {
			if util.Similarity(program.Name, game.Name) >= 0.6 {
				program.Name = game.Name
				break
			}
		}
	}

	return program
}
