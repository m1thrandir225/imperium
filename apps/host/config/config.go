package config

import "github.com/m1thrandir225/imperium/apps/host/internal/video"

type Config struct {
	videoConfig video.Config
}

func LoadConfig() Config {
	return Config{}
}
