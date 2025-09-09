package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func NormalizeName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

func Similarity(a, b string) float64 {
	na, nb := NormalizeName(a), NormalizeName(b)
	minLen := len(na)
	if len(nb) < minLen {
		minLen = len(nb)
	}
	match := 0
	for i := 0; i < minLen; i++ {
		if na[i] == nb[i] {
			match++
		} else {
			break
		}
	}
	return float64(match) / float64(len(na))
}

func ShortPath(p string) string {
	vol := filepath.VolumeName(p)
	rest := strings.TrimPrefix(p, vol)
	parts := strings.Split(rest, string(filepath.Separator))

	if len(parts) == 0 {
		return p
	}
	if len(parts) == 1 {
		if vol != "" {
			return filepath.Join(vol, parts[0])
		}
		return parts[0]
	}

	tail := filepath.Join(parts[len(parts)-2], parts[len(parts)-1])
	if vol != "" {
		return filepath.Join(vol, "...", tail)
	}
	return filepath.Join("...", tail)
}
