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
	"github.com/m1thrandir225/imperium/apps/host/pkg/rawg"
)

// Service represents the main functionality of the programs part
// Handles Program Discovery, Saving and Launching
// Rawg API Integration
type Service interface {
	DiscoverAndSavePrograms(paths []string) error
	GetLocalPrograms() ([]*Program, error)
	GetLocalProgramByPath(path string) (*Program, error)
	GetLocalProgramByID(id string) (*Program, error)
	DiscoverPrograms() ([]Program, error)
	DiscoverProgramsIn(paths []string) ([]Program, error) //??? should return pointer no?
	SaveProgram(req CreateProgramRequest) (*Program, error)
	LaunchProgram(path string) (*exec.Cmd, error)
	GetWindowTitleByProcessID(pid uint32) (string, error)
	RawgSearch(program Program) Program
}

type programService struct {
	db         Database
	rawgClient rawg.Client
}

// NewService creates a new program service
func NewService(
	dbPath,
	rawgAPIKey string,
) (Service, error) {
	return newProgramService(dbPath, rawgAPIKey)
}

func newProgramService(
	dbPath,
	rawgAPIKey string,
) (*programService, error) {
	db, err := NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	rawgClient, err := rawg.New(rawgAPIKey)
	if err != nil {
		return nil, err
	}

	return &programService{
		db:         db,
		rawgClient: rawgClient,
	}, nil
}

// DiscoverAndSavePrograms discovers program in common and custom paths and saves them to the local db
func (s *programService) DiscoverAndSavePrograms(paths []string) error {
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
func (s *programService) GetLocalPrograms() ([]*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}

	return s.db.GetPrograms()
}

// GetLocalProgramByPath gets a program from the local db by path
func (s *programService) GetLocalProgramByPath(path string) (*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}
	return s.db.GetProgramByPath(path)
}

func (s *programService) GetLocalProgramByID(id string) (*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}
	return s.db.GetProgramByID(id)
}

// DiscoverPrograms discovers programs in the common paths
func (s *programService) DiscoverPrograms() ([]Program, error) {
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
func (s *programService) DiscoverProgramsIn(paths []string) ([]Program, error) {
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

func (s *programService) scanDirectoryForPrograms(path string) ([]Program, error) {
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

func (s *programService) SaveProgram(req CreateProgramRequest) (*Program, error) {
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

func (s *programService) LaunchProgram(path string) (*exec.Cmd, error) {
	cmd := exec.Command(path)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to launch program: %w", err)
	}
	return cmd, nil
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

func (s *programService) RawgSearch(program Program) Program {
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
