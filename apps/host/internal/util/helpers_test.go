package util

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestCheckFFMPEGInstallation(t *testing.T) {}

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
	//TODO
}
