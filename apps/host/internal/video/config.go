// Package video provides a wrapper for the ffmpeg command line tool, and takes care of the encoding process.
package video

type Config struct {
	Encoder    string `json:"encoder" mapstructure:"encoder"`
	FPS        int    `json:"fps" mapstructure:"fps"`
	FFMPEGPath string `json:"ffmpeg_path" mapstructure:"ffmpeg_path"`
}

func (c *Config) SetEncoder(encoder string) {
	c.Encoder = encoder
}

func (c *Config) SetFPS(fps int) {
	c.FPS = fps
}

func (c *Config) SetFFMPEGPath(ffmpegPath string) {
	c.FFMPEGPath = ffmpegPath
}

func LoadConfig(encoder string, fps int, ffmpegPath string) *Config {
	return &Config{
		Encoder:    encoder,
		FPS:        fps,
		FFMPEGPath: ffmpegPath,
	}
}

// NewDefaultConfig returns a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Encoder:    "libx264",
		FPS:        30,
		FFMPEGPath: "",
	}
}
