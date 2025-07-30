package video

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Recorder struct {
	ffmpeg *FFMPEGWrapper
	config *Config
}

func NewRecorder(config *Config) *Recorder {
	return &Recorder{
		config: config,
		ffmpeg: NewFFMPEGWrapper(config.FFMPEGPath),
	}
}

func (r *Recorder) buildBaseArgs() []string {
	return []string{
		"-framerate", fmt.Sprintf("%d", r.config.FPS),
		"-probesize", "42M",
		"-thread_queue_size", "4096",
	}
}

func (r *Recorder) buildOutputArgs(outputPath string) []string {
	return []string{
		"-c:v", r.config.Encoder,
		"-preset", "ultrafast",
		"-pix-fmt", "yuv420p",
		"-y", // OVERWRITE IF EXISTS
		outputPath,
	}
}

// RecordScreen TODO: add linux/macos support for screen recording
func (r *Recorder) RecordScreen(outputPath string) error {
	args := r.buildBaseArgs()
	switch runtime.GOOS {
	case "windows":
		args = append(args, "-f", "gdigrab", "-i", "desktop")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	args = append(args, r.buildOutputArgs(outputPath)...)
	return r.ffmpeg.Execute(args...)
}

// RecordWindow TODO: add linux/macos support for window recording
func (r *Recorder) RecordWindow(windowTitle string, outputPath string) error {
	args := r.buildBaseArgs()

	switch runtime.GOOS {
	case "windows":
		args = append(args, "-f", "gdigrab", "-i", fmt.Sprintf("title=%s", windowTitle))
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	args = append(args, r.buildOutputArgs(outputPath)...)
	return r.ffmpeg.Execute(args...)
}

func (r *Recorder) StopRecording() error {
	return r.ffmpeg.Execute("-quit")
}

func (r *Recorder) GetTempOutputPath() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return filepath.Join(os.TempDir(), fmt.Sprintf("recording_%s.mp4", timestamp))
}
