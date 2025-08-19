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
	} else {
		return []string{
			"-preset", "llhq",
			"-bf", "0",
			"-b:v", "8M",
			"-maxrate", "10M",
			"-bufsize", "10M",
			"-spatial_aq", "1",
			"-rc-lookahead", "0",
			// repeat SPS/PPS and add AUD
			"-bsf:v", "dump_extra", "-bsf:v", "h264_metadata=aud=insert",
		}
	}
}

// RecordScreen TODO: add linux/macos support for screen recording
func (r *Recorder) RecordScreen(outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()
	switch runtime.GOOS {
	case "windows":
		args = append(args, "-f", "gdigrab", "-i", "desktop")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	args = append(args, r.convertRGBToBT709()...)
	args = append(args, r.buildCommonLowLatencyArgs()...)
	args = append(args, r.buildEncoderArgs(r.config.Encoder)...)
	args = append(args, r.buildStdOutputArgs()...)
	return r.ffmpeg.ExecuteWithStdout(args...)
}

// RecordWindow TODO: add linux/macos support for window recording
func (r *Recorder) RecordWindow(windowTitle string, outputPath *string) (io.ReadCloser, error) {
	args := r.buildBaseArgs()

	switch runtime.GOOS {
	case "windows":
		args = append(args, "-f", "gdigrab", "-i", fmt.Sprintf("title=%s", windowTitle))

	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	args = append(args, r.convertRGBToBT709()...)
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
