// Package programs provides the programs service for the host application.
package programs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
)

type ProgramService struct {
	authServerBaseURL string
	httpClient        *httpclient.Client
	db                *ProgramDB
}

func NewService(authServerBaseURL string,
	httpClient *httpclient.Client,
	dbPath string,
) *ProgramService {
	db, err := NewProgramDB(dbPath)
	if err != nil {
		//TODO: maybe we should panic
		log.Printf("Failed to initialize program database: %v", err)
	}
	return &ProgramService{
		authServerBaseURL: authServerBaseURL,
		httpClient:        httpClient,
		db:                db,
	}
}

func (s *ProgramService) DiscoverAndSavePrograms() error {
	programs, err := s.DiscoverPrograms()
	if err != nil {
		return err
	}

	for _, program := range programs {
		if s.db != nil {
			if err := s.db.SaveProgram(&program); err != nil {
				log.Printf("Failed to save program %s: %v", program.Name, err)
			}
		}
	}

	if s.db != nil {
		if err := s.db.CleanupNonExistentPrograms(); err != nil {
			log.Printf("Failed to cleanup non-existent programs: %v", err)
		}
	}

	return nil
}

func (s *ProgramService) GetLocalPrograms() ([]*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}

	return s.db.GetPrograms()
}

func (s *ProgramService) GetLocalProgramByPath(path string) (*Program, error) {
	if s.db == nil {
		return nil, fmt.Errorf("program database not initialized")
	}
	return s.db.GetProgramByPath(path)
}

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

func (s *ProgramService) scanDirectoryForPrograms(path string) ([]Program, error) {
	var programs []Program

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			discoveredPrograms, err := s.scanDirectoryForPrograms(filepath.Join(path, file.Name()))
			if err != nil {
				continue
			}
			programs = append(programs, discoveredPrograms...)
		}
	}

	return programs, nil
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

func (s *ProgramService) RegisterProgram(ctx context.Context, req CreateProgramRequest, hostID string) (*Program, error) {
	url := fmt.Sprintf("/api/v1/hosts/%s/programs", hostID)

	resp, err := s.httpClient.Post(ctx, url, req, make(map[string]string), true, make(map[string]string))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register program failed: %d", resp.StatusCode)
	}

	var result Program
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *ProgramService) GetProgramByID(ctx context.Context, programID string) (*Program, error) {
	url := fmt.Sprintf("/api/v1/programs/%s", programID)

	resp, err := s.httpClient.Get(ctx, url, make(map[string]string), make(map[string]string), true)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get program failed: %d", resp.StatusCode)
	}

	var program Program
	if err := json.Unmarshal(resp.Body, &program); err != nil {
		return nil, err
	}

	return &program, nil
}
