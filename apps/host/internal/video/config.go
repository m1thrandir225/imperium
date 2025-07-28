package video

type Config struct {
	Encoder    string `json:"encoder"`
	FPS        int    `json:"fps"`
	FFMPEGPath string `json:"ffmpeg_path"`
}

func LoadConfig(encoder string, fps int, ffmpegPath string) *Config {
	return &Config{
		Encoder:    encoder,
		FPS:        fps,
		FFMPEGPath: ffmpegPath,
	}
}
