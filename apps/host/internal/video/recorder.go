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

// Recorder takes care of recording video via ffmpeg
type Recorder struct {
	ffmpeg *FFMPEGWrapper
	config *Config
}

// NewRecorder returns a new Recorder instance based on the given config
func NewRecorder(config *Config) (*Recorder, error) {
	path := "ffmpeg"
	if config != nil && config.FFMPEGPath != "" {
		path = config.FFMPEGPath
	}

	ffmpegWrapper, err := NewFFMPEGWrapper(path)
	if err != nil {
		return nil, err
	}

	return &Recorder{
		config: config,
		ffmpeg: ffmpegWrapper,
	}, nil
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
		"-pix_fmt", "yuv420p", // This will be overridden for hardware encoders
		"-bf", "0",
		"-b_strategy", "0",
		"-sc_threshold", "0",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-g", fmt.Sprintf("%d", r.config.FPS), // ~1s GOP
		"-force_key_frames", "expr:gte(t,n_forced*1)",
		"-fps_mode", "cfr", // Use CFR instead of passthrough
		"-r", fmt.Sprintf("%d", r.config.FPS), // Set target framerate
	}
}

func (r *Recorder) buildEncoderArgs(encoder string) []string {
	if strings.Contains(encoder, "libx264") {
		return []string{
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-crf", "18", // Higher quality for games
			"-profile:v", "baseline",
			"-level:v", "3.1",
			"-x264-params", fmt.Sprintf("repeat-headers=1:scenecut=0:keyint=%d:min-keyint=%d:no-mbtree:no-cabac:no-deblock", r.config.FPS, r.config.FPS),
		}
	} else if strings.Contains(encoder, "nvenc") {
		return []string{
			"-preset", "p1",
			"-tune", "ll",
			"-bf", "0",
			"-rc", "cbr",
			"-b:v", "8M",
			"-maxrate", "8M",
			"-bufsize", "2M",
			"-surfaces", "1",
			"-spatial_aq", "0",
			"-rc-lookahead", "0",
			"-forced-idr", "1",
			"-no-scenecut", "1",
			"-delay", "0",
			"-bsf:v", "h264_metadata=aud=insert",
		}
	} else if strings.Contains(encoder, "libx265") {
		return []string{
			"-preset", "ultrafast",
			"-tune", "zerolatency",
			"-crf", "20",
			"-x265-params", "no-scenecut:keyint=60:min-keyint=60:no-mbtree:no-cabac:no-deblock",
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

// RecordScreen unified method - auto-selects capture based on encoder
func (r *Recorder) RecordScreen(outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()

	switch runtime.GOOS {
	case "windows":
		// Auto-select capture method based on encoder
		if strings.Contains(r.config.Encoder, "nvenc") || strings.Contains(r.config.Encoder, "qsv") || strings.Contains(r.config.Encoder, "amf") {

			w, h := 1920, 1080
			if mi, err := GetPrimaryMonitorInfo(); err == nil && mi.Width > 0 && mi.Height > 0 {
				w, h = mi.Width, mi.Height
				fmt.Printf("Using gfxcapture with primary monitor: %dx%d\n", w, h)
			}

			filter := fmt.Sprintf(
				"gfxcapture=monitor_idx=0:width=%d:height=%d:resize_mode=scale_aspect:output_fmt=8bit,hwdownload,format=bgra,format=nv12",
				w, h,
			)
			args = append(args, "-filter_complex", filter)
			args = append(args, "-filter_threads", "2")
			args = append(args, "-filter_complex_threads", "2")

		} else {

			monitorInfo, err := GetPrimaryMonitorInfo()
			if err != nil {
				fmt.Printf("Warning: Could not detect primary monitor, using full desktop: %v\n", err)
				args = append(args, "-f", "gdigrab")
				args = append(args, "-i", "desktop")
			} else {
				args = append(args, "-f", "gdigrab")
				args = append(args, "-i", "desktop")
				cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d",
					monitorInfo.Width,
					monitorInfo.Height,
					monitorInfo.OffsetX,
					monitorInfo.OffsetY)
				args = append(args, "-vf", cropFilter)
				fmt.Printf("Using gdigrab with primary monitor: %dx%d at offset (%d,%d)\n",
					monitorInfo.Width, monitorInfo.Height, monitorInfo.OffsetX, monitorInfo.OffsetY)
			}
		}

		args = append(args, "-fflags", "fastseek+flush_packets")
		args = append(args, "-flags", "global_header")

	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	args = append(args, r.buildCommonLowLatencyArgs()...)
	args = append(args, r.buildEncoderArgs(r.config.Encoder)...)

	if strings.Contains(r.config.Encoder, "nvenc") || strings.Contains(r.config.Encoder, "qsv") || strings.Contains(r.config.Encoder, "amf") {
		args = append(args, "-pix_fmt", "nv12")
	}

	args = append(args, r.buildStdOutputArgs()...)
	return r.ffmpeg.ExecuteWithStdout(args...)
}

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
