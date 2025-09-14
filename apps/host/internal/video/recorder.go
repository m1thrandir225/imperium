package video

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Recorder struct {
	ffmpeg *FFMPEGWrapper
	config *Config
}

func NewRecorder(config *Config) *Recorder {
	path := "ffmpeg"
	if config != nil && config.FFMPEGPath != "" {
		path = config.FFMPEGPath
	}
	return &Recorder{
		config: config,
		ffmpeg: NewFFMPEGWrapper(path),
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
		"-preset", "medium",
		"-pix_fmt", "yuv420p",
		"-y", // OVERWRITE IF EXISTS
		outputPath,
	}
}

func (r *Recorder) buildStdOutputArgs() []string {
	return []string{
		"-f", "h264",
		"-",
	}
}

func (r *Recorder) convertRGBToBT709() []string {
	return []string{
		"-vf", "colorspace=all=bt709:iall=bt709:range=tv,format=yuv420p",
		"-colorspace", "bt709", "-color_trc", "bt709", "-color_primaries", "bt709", "-color_range", "tv",
	}
}

func (r *Recorder) buildCommonLowLatencyArgs() []string {
	return []string{
		"-an",
		"-c:v", r.config.Encoder,
		"-pix_fmt", "yuv420p",
		"-bf", "0",
		"-b_strategy", "0",
		"-sc_threshold", "0",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-g", fmt.Sprintf("%d", r.config.FPS), // ~1s GOP
		"-force_key_frames", "expr:gte(t,n_forced*1)",
	}
}

func (r *Recorder) buildEncoderArgs(encoder string) []string {
	if strings.Contains(encoder, "libx264") {
		return []string{
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-crf", "20",
			"-profile:v", "baseline",
			"-level:v", "3.1",
			"-x264-params", fmt.Sprintf("repeat-headers=1:scenecut=0:keyint=%d:min-keyint=%d", r.config.FPS, r.config.FPS),
		}
	} else if strings.Contains(encoder, "nvenc") {
		return []string{
			"-preset", "p1", // Modern NVENC preset (fastest)
			"-tune", "ll", // Low latency tuning
			"-bf", "0",
			"-b:v", "8M",
			"-maxrate", "10M",
			"-bufsize", "10M",
			"-spatial_aq", "1",
			"-rc-lookahead", "0",
			// Only use one BSF to avoid conflicts
			"-bsf:v", "h264_metadata=aud=insert",
		}
	} else {
		return []string{
			"-preset", "medium",
			"-bf", "0",
			"-b:v", "8M",
			"-maxrate", "10M",
			"-bufsize", "10M",
		}
	}
}

// RecordWindow TODO: add linux/macos support for window recording
func (r *Recorder) RecordWindow(windowTitle string, outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()

	switch runtime.GOOS {
	case "windows":
		// Use GDI capture with optimizations for games
		args = append(args, "-f", "gdigrab")
		args = append(args, "-i", fmt.Sprintf("title=%s", windowTitle))
		args = append(args, "-draw_mouse", "0")
		args = append(args, "-show_region", "1")

		// Add some optimizations for better game capture
		args = append(args, "-fflags", "nobuffer+fastseek")
		args = append(args, "-flags", "low_delay")

	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	args = append(args, r.convertRGBToBT709()...)
	args = append(args, r.buildCommonLowLatencyArgs()...)
	args = append(args, r.buildEncoderArgs(r.config.Encoder)...)
	args = append(args, r.buildStdOutputArgs()...)
	return r.ffmpeg.ExecuteWithStdout(args...)
}

// RecordScreen TODO: add linux/macos support for screen recording
func (r *Recorder) RecordScreen(outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()
	switch runtime.GOOS {
	case "windows":
		// Get primary monitor info for multi-monitor setups
		monitorInfo, err := GetPrimaryMonitorInfo()
		if err != nil {
			fmt.Printf("Warning: Could not detect primary monitor, using full desktop: %v\n", err)
			// Fallback to full desktop capture
			args = append(args, "-f", "gdigrab")
			args = append(args, "-i", "desktop")
		} else {
			// Use primary monitor with cropping
			args = append(args, "-f", "gdigrab")
			args = append(args, "-i", "desktop")
			// Crop to primary monitor
			cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d",
				monitorInfo.Width,
				monitorInfo.Height,
				monitorInfo.OffsetX,
				monitorInfo.OffsetY)
			args = append(args, "-vf", cropFilter)
			fmt.Printf("Using primary monitor: %dx%d at offset (%d,%d)\n",
				monitorInfo.Width, monitorInfo.Height, monitorInfo.OffsetX, monitorInfo.OffsetY)
		}

		// Game-specific optimizations
		args = append(args, "-draw_mouse", "0")
		args = append(args, "-show_region", "1")
		args = append(args, "-fflags", "nobuffer+fastseek+flush_packets")
		args = append(args, "-flags", "low_delay+global_header")

	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Skip color conversion for games to reduce latency
	// args = append(args, r.convertRGBToBT709()...)

	args = append(args, r.buildCommonLowLatencyArgs()...)
	args = append(args, r.buildEncoderArgs(r.config.Encoder)...)
	args = append(args, r.buildStdOutputArgs()...)
	return r.ffmpeg.ExecuteWithStdout(args...)
}

// RecordGameScreen records the screen optimized specifically for games
// Uses primary monitor detection and game-optimized settings
func (r *Recorder) RecordGameScreen(outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()
	switch runtime.GOOS {
	case "windows":
		// Get primary monitor info for multi-monitor setups
		monitorInfo, err := GetPrimaryMonitorInfo()
		if err != nil {
			fmt.Printf("Warning: Could not detect primary monitor, using full desktop: %v\n", err)
			// Fallback to full desktop capture
			args = append(args, "-f", "gdigrab")
			args = append(args, "-i", "desktop")
		} else {
			// Use primary monitor with cropping
			args = append(args, "-f", "gdigrab")
			args = append(args, "-i", "desktop")
			// Crop to primary monitor
			cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d",
				monitorInfo.Width,
				monitorInfo.Height,
				monitorInfo.OffsetX,
				monitorInfo.OffsetY)
			args = append(args, "-vf", cropFilter)
			fmt.Printf("Game capture using primary monitor: %dx%d at offset (%d,%d)\n",
				monitorInfo.Width, monitorInfo.Height, monitorInfo.OffsetX, monitorInfo.OffsetY)
		}

		// Enhanced game-specific optimizations
		args = append(args, "-draw_mouse", "0")
		args = append(args, "-show_region", "1")
		args = append(args, "-fflags", "nobuffer+fastseek+flush_packets")
		args = append(args, "-flags", "low_delay+global_header")

		// Additional game optimizations
		args = append(args, "-avoid_negative_ts", "make_zero")
		args = append(args, "-fflags", "+genpts")

	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Skip color conversion for games to reduce latency
	args = append(args, r.buildCommonLowLatencyArgs()...)
	args = append(args, r.buildEncoderArgs(r.config.Encoder)...)
	args = append(args, r.buildStdOutputArgs()...)
	return r.ffmpeg.ExecuteWithStdout(args...)
}

// StopRecording stops the recording process
func (r *Recorder) StopRecording() error {
	if r.ffmpeg != nil {
		return r.ffmpeg.Stop()
	}
	return nil
}

// GetTempOutputPath returns a temporary output path for the recording
func (r *Recorder) GetTempOutputPath() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return filepath.Join(os.TempDir(), fmt.Sprintf("recording_%s.mp4", timestamp))
}

func (r *Recorder) GetFPS() int {
	if r.config != nil {
		return r.config.FPS
	}
	return 30
}
