//go:build darwin
// +build darwin

package video

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	cachedPrimaryMonitor *MonitorInfo
	lastPrimaryFetch     time.Time
)

// DisplayInfo represents display information from system_profiler
type DisplayInfo struct {
	Displays []GPUDevice `json:"SPDisplaysDataType"`
}

type GPUDevice struct {
	Name     string    `json:"_name"`
	Displays []Display `json:"spdisplays_ndrvs"`
}

type Display struct {
	Name       string `json:"_name"`
	Resolution string `json:"_spdisplays_resolution"`
	Pixels     string `json:"_spdisplays_pixels"`
	Main       string `json:"spdisplays_main,omitempty"`
	DisplayID  string `json:"_spdisplays_displayID,omitempty"`
}

// GetPrimaryMonitorInfo returns information about the primary monitor on macOS
func GetPrimaryMonitorInfo() (*MonitorInfo, error) {
	// Return cached result if it's recent
	if cachedPrimaryMonitor != nil && time.Since(lastPrimaryFetch) < 10*time.Minute {
		return cachedPrimaryMonitor, nil
	}

	monitors, err := GetAllMonitorsInfo()
	if err != nil {
		return nil, err
	}

	// Find the primary monitor
	for _, monitor := range monitors {
		if monitor.IsPrimary {
			cachedPrimaryMonitor = monitor
			lastPrimaryFetch = time.Now()
			return monitor, nil
		}
	}

	return nil, fmt.Errorf("primary monitor not found")
}

// GetMonitorCount returns the number of connected monitors on macOS
func GetMonitorCount() (int, error) {
	monitors, err := GetAllMonitorsInfo()
	if err != nil {
		return 0, err
	}
	return len(monitors), nil
}

// GetAllMonitorsInfo returns information about all connected monitors on macOS
func GetAllMonitorsInfo() ([]*MonitorInfo, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute system_profiler: %v", err)
	}

	var displayInfo DisplayInfo
	if err := json.Unmarshal(output, &displayInfo); err != nil {
		return nil, fmt.Errorf("failed to parse display info: %v", err)
	}

	var monitors []*MonitorInfo

	for _, gpu := range displayInfo.Displays {
		for _, display := range gpu.Displays {
			monitor, err := parseDisplayInfo(display)
			if err != nil {
				continue // Skip displays we can't parse
			}
			monitors = append(monitors, monitor)
		}
	}

	if len(monitors) == 0 {
		return getMonitorsFromDisplays()
	}

	return monitors, nil
}

// parseDisplayInfo converts Display struct to MonitorInfo
func parseDisplayInfo(display Display) (*MonitorInfo, error) {
	// Parse width and height
	width := 0
	height := 0

	if display.Pixels != "" {
		parts := strings.Split(display.Pixels, " x ")
		if len(parts) == 2 {
			var err error
			width, err = strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("failed to parse width from pixels: %v", err)
			}
			height, err = strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("failed to parse height from pixels: %v", err)
			}
		}
	} else if display.Resolution != "" {
		// Parse resolution string like "1920 x 1080 @ 165.00Hz"
		parts := strings.Split(display.Resolution, " @ ")
		if len(parts) >= 1 {
			resParts := strings.Split(parts[0], " x ")
			if len(resParts) == 2 {
				var err error
				width, err = strconv.Atoi(strings.TrimSpace(resParts[0]))
				if err != nil {
					return nil, fmt.Errorf("failed to parse width from resolution: %v", err)
				}
				height, err = strconv.Atoi(strings.TrimSpace(resParts[1]))
				if err != nil {
					return nil, fmt.Errorf("failed to parse height from resolution: %v", err)
				}
			}
		}
	}

	if width == 0 || height == 0 {
		return nil, fmt.Errorf("could not determine display dimensions")
	}

	isPrimary := display.Main == "spdisplays_yes"

	return &MonitorInfo{
		Width:     width,
		Height:    height,
		OffsetX:   0, // system_profiler doesn't provide offset info easily
		OffsetY:   0,
		IsPrimary: isPrimary,
	}, nil
}

// getMonitorsFromDisplays uses displayplacer as fallback to get monitor info
func getMonitorsFromDisplays() ([]*MonitorInfo, error) {
	cmd := exec.Command("displayplacer", "list")
	output, err := cmd.Output()
	if err != nil {
		return getBasicMonitorInfo()
	}

	return parseDisplayplacerOutput(string(output))
}

// getBasicMonitorInfo is the most basic fallback method
func getBasicMonitorInfo() ([]*MonitorInfo, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get basic display info: %v", err)
	}

	// Parse the text output for resolution information
	lines := strings.Split(string(output), "\n")
	var monitors []*MonitorInfo

	for i, line := range lines {
		if strings.Contains(line, "Resolution:") {
			// Extract resolution from line like "Resolution: 1920 x 1080"
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				resParts := strings.Split(strings.TrimSpace(parts[1]), " x ")
				if len(resParts) == 2 {
					width, err1 := strconv.Atoi(strings.TrimSpace(resParts[0]))
					height, err2 := strconv.Atoi(strings.TrimSpace(resParts[1]))

					if err1 == nil && err2 == nil {
						// Check if this is main display (look for "Main Display" in previous lines)
						isPrimary := false
						for j := max(0, i-10); j < i; j++ {
							if strings.Contains(strings.ToLower(lines[j]), "main") ||
								strings.Contains(strings.ToLower(lines[j]), "built-in") {
								isPrimary = true
								break
							}
						}

						monitors = append(monitors, &MonitorInfo{
							Width:     width,
							Height:    height,
							OffsetX:   0,
							OffsetY:   0,
							IsPrimary: isPrimary,
						})
					}
				}
			}
		}
	}

	if len(monitors) == 0 {
		return nil, fmt.Errorf("no monitors found")
	}

	foundPrimary := false
	for _, monitor := range monitors {
		if monitor.IsPrimary {
			foundPrimary = true
			break
		}
	}
	if !foundPrimary && len(monitors) > 0 {
		monitors[0].IsPrimary = true
	}

	return monitors, nil
}

// parseDisplayplacerOutput parses output from displayplacer tool
func parseDisplayplacerOutput(output string) ([]*MonitorInfo, error) {
	lines := strings.Split(output, "\n")
	var monitors []*MonitorInfo

	for _, line := range lines {
		if strings.Contains(line, "Resolution:") && strings.Contains(line, "Origin:") {
			monitor := parseDisplayplacerLine(line)
			if monitor != nil {
				monitors = append(monitors, monitor)
			}
		}
	}

	if len(monitors) == 0 {
		return nil, fmt.Errorf("no monitors found in displayplacer output")
	}

	// Mark the monitor at origin (0,0) as primary
	for _, monitor := range monitors {
		if monitor.OffsetX == 0 && monitor.OffsetY == 0 {
			monitor.IsPrimary = true
			break
		}
	}

	return monitors, nil
}

// parseDisplayplacerLine parses a single line from displayplacer output
func parseDisplayplacerLine(line string) *MonitorInfo {
	// Example line: "Resolution: 1920x1080 Origin: (0,0) Degree: 0"

	// Extract resolution
	resPart := extractBetween(line, "Resolution: ", " ")
	if resPart == "" {
		return nil
	}

	resParts := strings.Split(resPart, "x")
	if len(resParts) != 2 {
		return nil
	}

	width, err1 := strconv.Atoi(resParts[0])
	height, err2 := strconv.Atoi(resParts[1])
	if err1 != nil || err2 != nil {
		return nil
	}

	// Extract origin
	originPart := extractBetween(line, "Origin: (", ")")
	if originPart == "" {
		return &MonitorInfo{
			Width:     width,
			Height:    height,
			OffsetX:   0,
			OffsetY:   0,
			IsPrimary: false,
		}
	}

	originParts := strings.Split(originPart, ",")
	if len(originParts) != 2 {
		return nil
	}

	offsetX, err1 := strconv.Atoi(strings.TrimSpace(originParts[0]))
	offsetY, err2 := strconv.Atoi(strings.TrimSpace(originParts[1]))
	if err1 != nil || err2 != nil {
		offsetX, offsetY = 0, 0
	}

	return &MonitorInfo{
		Width:     width,
		Height:    height,
		OffsetX:   offsetX,
		OffsetY:   offsetY,
		IsPrimary: false, // Will be set later based on origin
	}
}

// extractBetween extracts substring between start and end markers
func extractBetween(str, start, end string) string {
	startIdx := strings.Index(str, start)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(start)

	endIdx := strings.Index(str[startIdx:], end)
	if endIdx == -1 {
		return str[startIdx:]
	}

	return str[startIdx : startIdx+endIdx]
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
