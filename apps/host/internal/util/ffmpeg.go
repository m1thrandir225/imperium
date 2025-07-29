package util

import (
	"os/exec"
	"runtime"
)

func CheckFFMPEGInstallation() (bool, string) {
	var paths []string
	switch runtime.GOOS {
	case "windows":
		paths = []string{"ffmpeg.exe", "C:\\ffmpeg\\bin\\ffmpeg.exe"}
	default:
		paths = []string{"/usr/bin/ffmpeg", "/usr/local/bin/ffmpeg"}
	}

	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		path, _ := exec.LookPath("ffmpeg")
		return true, path
	}

	for _, path := range paths {
		cmd := exec.Command(path, "-version")
		if err := cmd.Run(); err != nil {
			return true, path
		}
	}
	return false, ""
}
