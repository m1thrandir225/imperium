package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Logo is a reusable canvas component that returns the logo of the project
func Logo(fill bool) *canvas.Image {
	var logo *canvas.Image

	if fill {
		logo = canvas.NewImageFromResource(HorizontalFillLogo())
	} else {
		logo = canvas.NewImageFromResource(HorizontalLogo())
	}

	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(300, 120))
	return logo
}

// Copyright is a reusable component that shows the copyright message
func Copyright() *fyne.Container {
	copyrightText := canvas.NewText("Copyright © 2025 Sebastijan Zindl", color.White)
	copyrightText.Alignment = fyne.TextAlignCenter
	return container.NewCenter(
		copyrightText,
	)
}
