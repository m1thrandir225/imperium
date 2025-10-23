package resources

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	LogoURL = "https://github.com/m1thrandir225/imperium/blob/master/assets/imperium_horizontal_logo.png?raw=true"
)

// LogoManager handles logo loading and fallback
type LogoManager struct {
	cachedResource fyne.Resource
	loadError      error
}

var globalLogoManager *LogoManager

// GetLogoManager returns the singleton logo manager
func GetLogoManager() *LogoManager {
	if globalLogoManager == nil {
		globalLogoManager = &LogoManager{}
	}
	return globalLogoManager
}

// LoadLogo attempts to load the logo from URL and caches it
func (lm *LogoManager) LoadLogo() error {
	if lm.cachedResource != nil {
		return lm.loadError
	}

	resource, err := fyne.LoadResourceFromURLString(LogoURL)
	lm.cachedResource = resource
	lm.loadError = err
	return err
}

// CreateLogoImage creates a logo image with proper sizing and fallback
func (lm *LogoManager) CreateLogoImage(width, height float32) fyne.CanvasObject {
	// Ensure logo is loaded
	lm.LoadLogo()

	if lm.loadError != nil || lm.cachedResource == nil {
		// Create fallback text logo
		return lm.createFallbackLogo(width, height)
	}

	// Create image from resource
	logo := canvas.NewImageFromResource(lm.cachedResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(width, height))

	return container.NewCenter(logo)
}

// createFallbackLogo creates a text-based logo when image loading fails
func (lm *LogoManager) createFallbackLogo(width, height float32) fyne.CanvasObject {
	textLogo := widget.NewLabelWithStyle("IMPERIUM", fyne.TextAlignCenter, fyne.TextStyle{
		Bold: true,
	})

	separator := widget.NewSeparator()

	fallbackContainer := container.NewVBox(
		textLogo,
		separator,
	)

	centeredContainer := container.NewCenter(fallbackContainer)
	centeredContainer.Resize(fyne.NewSize(width, height))

	return centeredContainer
}

// PreloadResources preloads all resources at application startup
func PreloadResources() {
	logoManager := GetLogoManager()
	logoManager.LoadLogo()
}

// CreateStandardLogo creates a logo with standard dimensions (200x80)
func CreateStandardLogo() fyne.CanvasObject {
	return GetLogoManager().CreateLogoImage(200, 80)
}

// CreateLargeLogo creates a logo with larger dimensions (300x120)
func CreateLargeLogo() fyne.CanvasObject {
	return GetLogoManager().CreateLogoImage(300, 120)
}

// CreateSmallLogo creates a logo with smaller dimensions (150x60)
func CreateSmallLogo() fyne.CanvasObject {
	return GetLogoManager().CreateLogoImage(150, 60)
}
