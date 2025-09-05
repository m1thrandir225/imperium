package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
)

type ProgramsScreen struct {
	manager      *Manager
	programsList *widget.List
	programs     []programs.Program
}

func NewProgramsScreen(manager *Manager) *ProgramsScreen {
	return &ProgramsScreen{
		manager: manager,
	}
}

func (s *ProgramsScreen) Name() string {
	return PROGRAMS_SCREEN
}

func (s *ProgramsScreen) Render(w fyne.Window) fyne.CanvasObject {
	// Discover programs
	discoveredPrograms, err := s.programService.DiscoverPrograms()
	if err != nil {
		dialog.ShowError(err, w)
	} else {
		s.programs = discoveredPrograms
	}

	// Create list widget
	s.programsList = widget.NewList(
		func() int { return len(s.programs) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Program Name"),
				widget.NewLabel("Path"),
				widget.NewButton("Register", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			program := s.programs[id]

			nameLabel := box.Objects[0].(*widget.Label)
			pathLabel := box.Objects[1].(*widget.Label)
			registerBtn := box.Objects[2].(*widget.Button)

			nameLabel.SetText(program.Name)
			pathLabel.SetText(program.Path)
			registerBtn.OnTapped = func() {
				s.registerProgram(program, w)
			}
		},
	)

	refreshBtn := widget.NewButton("Refresh Programs", func() {
		discoveredPrograms, err := s.programService.DiscoverPrograms()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		s.programs = discoveredPrograms
		s.programsList.Refresh()
	})

	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	content := container.NewBorder(
		container.NewHBox(refreshBtn, backBtn),
		nil, nil, nil,
		s.programsList,
	)

	return content
}

func (s *ProgramsScreen) registerProgram(program programs.Program, w fyne.Window) {
	req := programs.CreateProgramRequest{
		Name:        program.Name,
		Path:        program.Path,
		Description: program.Description,
	}

	_, err := s.programService.RegisterProgram(context.Background(), req, "current-host-id")
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	dialog.ShowInformation("Success", "Program registered successfully!", w)
}
