package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AuthServerBaseURL string
	LoginPath         string
	HostsPath         string
	MatchmakePath     string

	CookieDomain string
	CookieSecure bool
}

func Load() *Config {
	authServerBaseURL := getenv("AUTH_SERVER_BASE_URL", "http://localhost:8080/api/v1")

	cookieDomain := getenv("COOKIE_DOMAIN", "")
	cookieSecure := getbool("COOKIE_SECURE", false)

	return &Config{
		AuthServerBaseURL: authServerBaseURL,
		LoginPath:         fmt.Sprintf("%s/login", authServerBaseURL),
		HostsPath:         fmt.Sprintf("%s/hosts", authServerBaseURL),
		MatchmakePath:     fmt.Sprintf("%s/matchmake", authServerBaseURL),
		CookieDomain:      cookieDomain,
		CookieSecure:      cookieSecure,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getbool(k string, def bool) bool {
	if v := os.Getenv(k); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}
