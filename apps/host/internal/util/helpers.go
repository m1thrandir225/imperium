package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GetConfigDir(appName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Failed to get user's home directory: %v", err)
		return "", err
	}

	return filepath.Join(configDir, appName), nil
}

func CheckFFMPEGInstallation() (bool, string) {
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		cmd := exec.Command(path, "-version")
		if err := cmd.Run(); err == nil {
			return true, path
		}
	}

	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		path, _ := exec.LookPath("ffmpeg")
		return true, path
	}

	var paths []string
	switch runtime.GOOS {
	case "windows":
		paths = []string{"ffmpeg.exe", "C:\\ffmpeg\\bin\\ffmpeg.exe"}
	default:
		paths = []string{"/usr/bin/ffmpeg", "/usr/local/bin/ffmpeg"}
	}

	for _, path := range paths {
		cmd := exec.Command(path, "-version")
		if err := cmd.Run(); err == nil {
			return true, path
		}
	}

	return false, ""
}
