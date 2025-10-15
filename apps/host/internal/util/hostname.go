package util

import "os"

// GetHostname returns the hostname of the system.
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}
