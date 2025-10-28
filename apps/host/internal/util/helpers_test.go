package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfigDir(t *testing.T) {
	appName := "test"

	configDir, err := os.UserConfigDir()
	require.NoError(t, err)
	require.NotEmpty(t, configDir)

	expectedPath := filepath.Join(configDir, appName)

	actualPath, err := GetConfigDir(appName)
	require.NoError(t, err)

	require.Equal(t, expectedPath, actualPath)
}

func TestCheckFFMPEGInstallation(t *testing.T) {
	installed, path := CheckFFMPEGInstallation()

	if installed {
		require.NotEmpty(t, path, "If FFMPEG is installed, path should not be empty")

		if path != "ffmpeg" && !filepath.IsAbs(path) {
			_, err := os.Stat(path)
			require.NoError(t, err, "Reported FFMPEG path should exist")
		}
	} else {
		require.Empty(t, path, "If FFMPEG is not installed, path should be empty")
	}

	require.IsType(t, true, installed)
	require.IsType(t, "", path)
}

func TestNormalizeName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Test ",
			expected: "Test",
		},
		{
			input:    "Test-Case",
			expected: "TestCase",
		},
		{
			input:    "Test - Case",
			expected: "TestCase",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, NormalizeName(tc.input))
	}
}

func TestSimilarity(t *testing.T) {
	testCases := []struct {
		inputA   string
		inputB   string
		expected float64
	}{
		{
			inputA:   "Test",
			inputB:   "Test",
			expected: 1.0,
		},
		{
			inputA:   "tast",
			inputB:   "Test",
			expected: 0.25,
		},
		{
			inputA:   "random word",
			inputB:   "hello world",
			expected: 0.0,
		},
	}

	for _, tc := range testCases {

		fmt.Println(Similarity(tc.inputA, tc.inputA))
		require.Equal(t, tc.expected, Similarity(tc.inputA, tc.inputB))
	}
}

func TestShortPath(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		name     string
	}{
		{
			name:     "single file",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "two level path",
			input:    filepath.Join("dir", "file.txt"),
			expected: filepath.Join("...", "dir", "file.txt"),
		},
		{
			name:     "three level path",
			input:    filepath.Join("path", "to", "file.txt"),
			expected: filepath.Join("...", "to", "file.txt"),
		},
		{
			name:     "deep path",
			input:    filepath.Join("very", "long", "path", "to", "file.txt"),
			expected: filepath.Join("...", "to", "file.txt"),
		},
		{
			name:     "root path single file",
			input:    filepath.Join(string(filepath.Separator), "file.txt"),
			expected: filepath.Join("...", "file.txt"),
		},
		{
			name:     "root path with subdirs",
			input:    filepath.Join(string(filepath.Separator), "usr", "bin", "ffmpeg"),
			expected: filepath.Join("...", "bin", "ffmpeg"),
		},
	}

	// Windows-specific test
	if runtime.GOOS == "windows" {
		windowsTests := []struct {
			input    string
			expected string
			name     string
		}{
			{
				name:     "windows drive single file",
				input:    "C:\\file.txt",
				expected: "C:\\file.txt",
			},
			{
				name:     "windows drive with path",
				input:    "C:\\Users\\John\\Documents\\file.txt",
				expected: "C:\\...\\Documents\\file.txt",
			},
			{
				name:     "windows drive two levels",
				input:    "D:\\Projects\\myapp",
				expected: "D:\\Projects\\myapp",
			},
		}
		testCases = append(testCases, windowsTests...)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ShortPath(tc.input)
			require.Equal(t, tc.expected, result, "ShortPath(%q) should return %q, got %q", tc.input, tc.expected, result)
		})
	}
}
