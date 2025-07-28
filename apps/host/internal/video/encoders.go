package video

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var h264Preferred = []string{
	"libx264",
	"h264_nvenc",
	"h264_qsv",
	"h264_amf",
	"h264_videotoolbox",
}

var h265Preferred = []string{
	"hevc_nvenc",
	"hevc_qsv",
	"hevc_amf",
	"hevc_videotoolbox",
	"libx265",
}

func itemInList(el string, list []string) bool {
	for _, value := range list {
		fmt.Println(value, el)
		if strings.Contains(value, el) {
			return true
		}
	}
	return false
}

func GetAvailableEncodersForCodecs() (h264Encoders, h265Encoders []string, err error) {
	cmd := exec.Command("ffmpeg", "-encoders")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		log.Printf("warning: 'ffmpeg -encoders' command finished with error: %v. Output may still be parsable.", err)
	}

	scanner := bufio.NewScanner(&out)
	inEncodersSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !inEncodersSection {
			if line == "Encoders:" {
				inEncodersSection = true
			}
			continue
		}
		if len(line) == 0 || line[0] != 'V' || strings.Contains(line, "=") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		encoderName := fields[1]

		if itemInList(encoderName, h264Preferred) {
			h264Encoders = append(h264Encoders, encoderName)
		} else if itemInList(encoderName, h265Preferred) {
			h265Encoders = append(h265Encoders, encoderName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading ffmpeg output: %w", err)
	}

	if !inEncodersSection {
		return nil, nil, fmt.Errorf("could not find 'Encoders:' section in ffmpeg output")
	}

	return h264Encoders, h265Encoders, nil
}
