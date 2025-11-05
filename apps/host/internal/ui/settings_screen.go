package ui

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
)

type SettingsScreen struct {
	manager *uiManager
}

func NewSettingsScreen(manager *uiManager) *SettingsScreen {
	return &SettingsScreen{
		manager: manager,
	}
}

func (s *SettingsScreen) Name() string {
	return SETTINGS_SCREEN
}

func (s *SettingsScreen) Render(w fyne.Window) fyne.CanvasObject {
	current := s.manager.GetState().Settings

	// Server Address Section
	serverAddressEntry := widget.NewEntry()
	serverAddressEntry.SetPlaceHolder("e.g., http://localhost:8080 or https://auth.example.com")
	if current.ServerAddress != "" {
		serverAddressEntry.SetText(current.ServerAddress)
	}

	// FFmpeg Path Section
	ffmpegPathEntry := widget.NewEntry()
	ffmpegPathEntry.SetPlaceHolder("Path to FFmpeg executable")
	if current.FFmpegPath != "" {
		ffmpegPathEntry.SetText(current.FFmpegPath)
	}

	browseFFmpegBtn := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			defer func() {
				if uri != nil {
					_ = uri.Close()
				}
			}()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			ffmpegPathEntry.SetText(uri.URI().Path())
		}, w)
	})

	// Encoder Selection - Start with basic fallbacks
	fallbackEncoders := []string{"libx264", "libx265"}
	encoderSelect := widget.NewSelect(fallbackEncoders, nil)
	if current.Encoder != "" {
		encoderSelect.SetSelected(current.Encoder)
	} else {
		encoderSelect.SetSelected("libx264") // Default
	}

	// Function to load available encoders
	loadAvailableEncoders := func() {
		h264Encoders, h265Encoders, err := video.GetAvailableEncodersForCodecs()
		if err != nil {
			// UI updates must be on main thread
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("failed to detect encoders: %v", err), w)
			})
			return
		}

		availableEncoders := make([]string, 0)
		availableEncoders = append(availableEncoders, h264Encoders...)
		availableEncoders = append(availableEncoders, h265Encoders...)

		fyne.Do(func() {
			if len(availableEncoders) > 0 {
				currentSelection := encoderSelect.Selected

				encoderSelect.Options = availableEncoders

				found := false
				if slices.Contains(availableEncoders, currentSelection) {
					encoderSelect.SetSelected(currentSelection)
					found = true
				}

				if !found && len(availableEncoders) > 0 {
					encoderSelect.SetSelected(availableEncoders[0])
				}

				encoderSelect.Refresh()
				dialog.ShowInformation("Success",
					fmt.Sprintf("Detected %d H.264 encoders and %d H.265 encoders",
						len(h264Encoders), len(h265Encoders)), w)
			} else {
				encoderSelect.Options = fallbackEncoders
				encoderSelect.SetSelected("libx264")
				encoderSelect.Refresh()
				dialog.ShowInformation("No Hardware Encoders",
					"No hardware encoders detected. Using software fallbacks.", w)
			}
		})
	}

	// Auto-load encoders on startup if FFmpeg path is available
	if current.FFmpegPath != "" {
		go loadAvailableEncoders() // Run in background to avoid blocking UI
	}

	// Detect encoders button
	detectEncodersBtn := widget.NewButton("Detect Available Encoders", func() {
		loadAvailableEncoders()
	})

	fpsOptions := []string{"30", "60", "90", "120"}
	fpsSelect := widget.NewSelect(fpsOptions, nil)
	fpsSelect.SetSelected(fmt.Sprintf("%d", current.Framerate))

	// Validation functions
	validateServerAddress := func(address string) error {
		if address == "" {
			return fmt.Errorf("server address is required")
		}

		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			return fmt.Errorf("server address must start with http:// or https://")
		}

		_, err := url.Parse(address)
		if err != nil {
			return fmt.Errorf("invalid server address: %v", err)
		}

		return nil
	}

	validateFFmpegPath := func(path string) error {
		if path == "" {
			return fmt.Errorf("FFmpeg path is required")
		}
		return nil
	}

	// Save button
	saveBtn := widget.NewButton("Save Settings", func() {
		if err := validateServerAddress(serverAddressEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		if err := validateFFmpegPath(ffmpegPathEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		fps, err := strconv.Atoi(fpsSelect.Selected)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid FPS value"), w)
			return
		}

		s.manager.publish(uapp.EventSettingsSaved, uapp.SettingsSavedPayload{
			Settings: state.Settings{
				FFmpegPath:    ffmpegPathEntry.Text,
				ServerAddress: serverAddressEntry.Text,
				Encoder:       encoderSelect.Selected,
				Framerate:     fps,
			},
		})

		dialog.ShowInformation("Success", "Settings saved successfully!", w)
	})

	// Reset button
	resetBtn := widget.NewButton("Reset to Defaults", func() {
		dialog.ShowConfirm("Reset Settings", "Are you sure you want to reset all settings to defaults?", func(confirmed bool) {
			if confirmed {
				serverAddressEntry.SetText("")
				ffmpegPathEntry.SetText("")
				encoderSelect.Options = fallbackEncoders
				encoderSelect.SetSelected("libx264")
				fpsSelect.SetSelected("30")
				encoderSelect.Refresh()
			}
		}, w)
	})

	// Test FFmpeg button
	testFFmpegBtn := widget.NewButton("Test FFmpeg", func() {
		if ffmpegPathEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please set FFmpeg path first"), w)
			return
		}

		// Test FFmpeg installation
		wrapper, err := video.NewFFMPEGWrapper(ffmpegPathEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("something went wrong: %v", err), w)
			return
		}
		version, err := wrapper.Version()
		if err != nil {
			dialog.ShowError(fmt.Errorf("FFmpeg test failed: %v", err), w)
			return
		}

		// Show first few lines of version info
		versionStr := string(version)
		lines := strings.Split(versionStr, "\n")
		if len(lines) > 3 {
			versionStr = strings.Join(lines[:3], "\n")
		}

		dialog.ShowInformation("FFmpeg Test", fmt.Sprintf("FFmpeg is working!\n\n%s", versionStr), w)
	})

	// Back button
	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	// Create form layout
	form := container.NewVBox(
		widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),

		// Server Section
		widget.NewLabel("Server Configuration"),
		widget.NewLabel("Server Address:"),
		serverAddressEntry,
		widget.NewSeparator(),

		// Video Section
		widget.NewLabel("Video Configuration"),
		widget.NewLabel("FFmpeg Path:"),
		container.NewBorder(nil, nil, nil, browseFFmpegBtn, ffmpegPathEntry),
		container.NewHBox(testFFmpegBtn, detectEncodersBtn),

		widget.NewLabel("Video Encoder:"),
		encoderSelect,

		widget.NewLabel("FPS (Frames Per Second):"),
		fpsSelect,
		widget.NewSeparator(),

		// Action buttons
		container.NewHBox(saveBtn, resetBtn),
		backBtn,
	)

	return container.NewScroll(container.NewPadded(form))
}
