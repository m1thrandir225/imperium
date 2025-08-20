package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

func GetUniquePCID() (string, error) {
	interfaces, err := net.Interfaces()

	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var macs []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if len(iface.HardwareAddr) > 0 {
				macs = append(macs, iface.HardwareAddr.String())
			}
		}
	}

	if len(macs) == 0 {
		return "", fmt.Errorf("could not find a valid MAC address")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	combinedMacs := strings.Join(macs, "-")

	combined := combinedMacs + "-" + hostname
	hash := sha256.Sum256([]byte(combined))

	return hex.EncodeToString(hash[:]), nil
}
