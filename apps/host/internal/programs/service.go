// Package programs provides the programs service for the host application.
package programs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ProgramService struct {
	authServerBaseURL string
	httpClient        *http.Client
	token             string
}

func NewProgramService(authServerBaseURL string, token string) *ProgramService {
	return &ProgramService{
		authServerBaseURL: authServerBaseURL,
		httpClient:        &http.Client{},
		token:             token,
	}
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
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/hosts/%s/programs", s.authServerBaseURL, hostID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register program failed: %s", resp.Status)
	}

	var result Program
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *ProgramService) GetProgramByID(ctx context.Context, programID string) (*Program, error) {
	url := fmt.Sprintf("%s/api/v1/programs/%s", s.authServerBaseURL, programID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get program failed: %s", resp.Status)
	}

	var program Program
	if err := json.NewDecoder(resp.Body).Decode(&program); err != nil {
		return nil, err
	}

	return &program, nil
}
