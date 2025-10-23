package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/imperium_horizontal_fill_logo.png
var horizontalFillLogo []byte

var horizontalFillLogoRes = &fyne.StaticResource{
	StaticName:    "imperium_horizontal_fill_logo.png",
	StaticContent: horizontalFillLogo,
}

//go:embed icons/imperium_horizontal_logo.png
var horizontalLogo []byte

var horizontalLogoRes = &fyne.StaticResource{
	StaticName:    "imperium_horizontal_logo.png",
	StaticContent: horizontalLogo,
}

//go:embed icons/imperium_icon_logo.png
var iconLogo []byte

var iconLogoRes = &fyne.StaticResource{
	StaticName:    "imperium_icon_logo.png",
	StaticContent: iconLogo,
}

//go:embed icons/imperium_vertical_logo.png
var verticalLogo []byte

var verticalLogoRes = &fyne.StaticResource{
	StaticName:    "imperium_vertical_logo.png",
	StaticContent: verticalLogo,
}
