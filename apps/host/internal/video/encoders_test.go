package video

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItemInList(t *testing.T) {
	cases := []struct {
		item   string
		list   []string
		result bool
	}{
		{
			item:   "ok",
			list:   []string{"no", "yes", "ok"},
			result: true,
		},
		{
			item:   "ok",
			list:   []string{"ok"},
			result: true,
		},
		{
			item:   "ok",
			list:   []string{"no", "yes"},
			result: false,
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.result, itemInList(tc.item, tc.list))
	}
}

func TestParseEncodersOutput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantH264    []string
		wantH265    []string
		expectError bool
	}{
		{
			name: "valid ffmpeg output",
			input: `
Encoders:
 V..... libx264             H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
 V..... libx265             H.265 / HEVC (codec hevc)
 V..... h264_nvenc          NVIDIA NVENC H.264 encoder
 V..... hevc_nvenc          NVIDIA NVENC HEVC encoder
 A..... libmp3lame          MP3 (MPEG audio layer 3)
`,
			wantH264:    []string{"libx264", "h264_nvenc"},
			wantH265:    []string{"libx265", "hevc_nvenc"},
			expectError: false,
		},
		{
			name: "missing Encoders section",
			input: `
Some other section:
 V..... libx264             H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
`,
			expectError: true,
		},
		{
			name:        "empty output",
			input:       ``,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h264, h265, err := parseEncodersOutput(tt.input)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, h264, tt.wantH264)
			require.Equal(t, h265, tt.wantH265)
		})
	}
}
