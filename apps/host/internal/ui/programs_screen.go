package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
)

type ProgramsScreen struct {
	manager      *Manager
	programsList *widget.List
	programs     []uapp.ProgramItem
	subscribed   bool
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
	if !s.subscribed {
		ch := s.manager.bus.Subscribe(uapp.EventProgramsDisocvered)
		s.subscribed = true

		go func() {
			for evt := range ch {
				payload, ok := evt.(uapp.ProgramDiscoveredPayload)
				if !ok {
					continue
				}
				s.programs = payload.Programs
				if s.programsList != nil {
					fyne.Do(func() { s.programsList.Refresh() })
				}
			}
		}()
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
				s.manager.Publish(uapp.EventProgramRegisterRequested, uapp.ProgramRegisterRequestedPayload{
					Program: program,
				})
			}
		},
	)

	refreshBtn := widget.NewButton("Refresh Programs", func() {
		s.manager.Publish(uapp.EventProgramsDiscoverRequested, nil)
	})

	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	//Request initial discovery
	s.manager.Publish(uapp.EventProgramsDiscoverRequested, nil)

	content := container.NewBorder(
		container.NewHBox(refreshBtn, backBtn),
		nil, nil, nil,
		s.programsList,
	)

	return content
}
