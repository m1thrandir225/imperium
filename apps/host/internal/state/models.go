package state

import "time"

type AppState struct {
	UserSession UserSession `mapstructure:"user_session" json:"user_session" yaml:"user_session"`
	UserInfo    UserInfo    `mapstructure:"user_info" json:"user_info" yaml:"user_info"`
	HostInfo    HostInfo    `mapstructure:"host_info" json:"host_info" yaml:"host_info"`
	Settings    Settings    `mapstructure:"settings" json:"settings" yaml:"settings"`
}

type UserSession struct {
	AccessToken           string    `mapstructure:"access_token" json:"access_token" yaml:"access_token"`
	RefreshToken          string    `mapstructure:"refresh_token" json:"refresh_token" yaml:"refresh_token"`
	AccessTokenExpiresAt  time.Time `mapstructure:"access_token_expires_at" json:"access_token_expires_at" yaml:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `mapstructure:"refresh_token_expires_at" json:"refresh_token_expires_at" yaml:"refresh_token_expires_at"`
}

type UserInfo struct {
	ID    string `mapstructure:"id" json:"id" yaml:"id"`
	Name  string `mapstructure:"name" json:"name" yaml:"name"`
	Email string `mapstructure:"email" json:"email" yaml:"email"`
}

type HostInfo struct {
	ID   string `mapstructure:"id" json:"id" yaml:"id"`       // unique host id from auth-server
	Name string `mapstructure:"name" json:"name" yaml:"name"` // friendly name for this host
	IP   string `mapstructure:"ip" json:"ip" yaml:"ip"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type Settings struct {
	FFmpegPath         string   `mapstructure:"ffmpeg_path" json:"ffmpeg_path" yaml:"ffmpeg_path"`
	Encoder            string   `mapstructure:"encoder" json:"encoder" yaml:"encoder"`
	Framerate          int      `mapstructure:"framerate" json:"framerate" yaml:"framerate"`
	Bitrate            string   `mapstructure:"bitrate" json:"bitrate" yaml:"bitrate"`
	ServerAddress      string   `mapstructure:"server_address" json:"server_address" yaml:"server_address"`
	CustomProgramPaths []string `mapstructure:"custom_program_paths" json:"custom_program_paths" yaml:"custom_program_paths"`
	RawgAPIKey         string   `mapstructure:"rawg_api_key" json:"rawg_api_key" yaml:"rawg_api_key"`
}
