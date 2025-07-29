package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"image/color"
	"net/url"
	"strings"
)

const SETUP_SCREEN = "Setup"

type SetupScreen struct {
	cfg        *util.Config
	saveConfig func(config *util.Config) error
	onComplete func()
}

func NewSetupScreen(cfg *util.Config, saveConfig func(config *util.Config) error, onComplete func()) *SetupScreen {
	return &SetupScreen{
		cfg:        cfg,
		saveConfig: saveConfig,
		onComplete: onComplete,
	}
}

func (s *SetupScreen) Name() string {
	return SETUP_SCREEN
}

func (s *SetupScreen) Render(w fyne.Window) fyne.CanvasObject {
	ffmpegStatusLabel := widget.NewLabel("Checking FFmpeg installation...")

	//ffmpeg path entry
	ffmpegPathEntry := widget.NewEntry()
	ffmpegPathEntry.Hide()

	//server address entry
	serverAddressEntry := widget.NewEntry()
	serverAddressEntry.SetPlaceHolder("e.g., http://localhost:8080 or https://auth.example.com")
	if s.cfg.ServerAddress != "" {
		serverAddressEntry.SetText(s.cfg.ServerAddress)
	}

	validateServerAddress := func(address string) error {
		if address == "" {
			return fmt.Errorf("server address is required")
		}

		// Basic URL validation
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			return fmt.Errorf("server address must start with http:// or https://")
		}

		_, err := url.Parse(address)
		if err != nil {
			return fmt.Errorf("invalid server address: %v", err)
		}

		return nil
	}

	browseBtn := widget.NewButton("Browse for FFmpeg", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			ffmpegPathEntry.SetText(uri.URI().Path())
			s.cfg.VideoConfig.FFMPEGPath = uri.URI().Path()
		}, w)
	})
	browseBtn.Hide()

	downloadURL, _ := url.Parse("https://ffmpeg.org/download.html")
	downloadLink := widget.NewHyperlink("Download FFmpeg", downloadURL)
	downloadLink.Hide()

	continueBtn := widget.NewButton("Continue", func() {
		// Validate server address
		if err := validateServerAddress(serverAddressEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		if ffmpegPathEntry.Text != "" {
			s.cfg.VideoConfig.FFMPEGPath = ffmpegPathEntry.Text
		}

		s.cfg.ServerAddress = serverAddressEntry.Text

		if err := s.saveConfig(s.cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}
		s.onComplete()
	})
	continueBtn.Hide()

	refreshBtn := widget.NewButton("Refresh", nil)
	refreshBtn.Hide()

	checkFFmpeg := func() {
		installed, path := util.CheckFFMPEGInstallation()
		if installed {
			ffmpegStatusLabel.SetText("✅ FFmpeg is installed")
			s.cfg.VideoConfig.FFMPEGPath = path
			continueBtn.Show()
			browseBtn.Hide()
			downloadLink.Hide()
			ffmpegPathEntry.Hide()
		} else {
			ffmpegStatusLabel.SetText("❌ FFmpeg is not installed")
			browseBtn.Show()
			downloadLink.Show()
			ffmpegPathEntry.Show()
			continueBtn.Show()
		}
		refreshBtn.Hide()
	}

	refreshBtn.OnTapped = checkFFmpeg

	checkFFmpeg()

	content := container.NewVBox(
		widget.NewLabel("Welcome to Imperium"),
		widget.NewLabel("Initial Setup"),
		canvas.NewLine(color.White),
		ffmpegStatusLabel,
		ffmpegPathEntry,
		container.NewHBox(browseBtn, refreshBtn),
		downloadLink,
		widget.NewLabel("Authentication Server Address"),
		serverAddressEntry,
		continueBtn,
	)

	return container.NewPadded(content)

}
