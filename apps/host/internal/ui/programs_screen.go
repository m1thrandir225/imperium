package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
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

	addProgramBtn := widget.NewButton("Add Program", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			defer func() {
				if uri != nil {
					_ = uri.Close()
				}
			}()
			if err != nil {
				dialog.ShowError(err, s.manager.window)
				return
			}
			if uri == nil {
				return
			}

			path := uri.URI().Path()
			name := filepath.Base(path)

			s.manager.Publish(uapp.EventProgramRegisterRequested, uapp.ProgramRegisterRequestedPayload{
				Program: uapp.ProgramItem{
					Name:        name,
					Path:        path,
					Description: "",
				},
			})
		}, w)
	})

	addScanPathBtn := widget.NewButton("Add Scan Path", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}

			current := s.manager.GetState().Settings.CustomProgramPaths
			newPath := uri.Path()

			found := false
			for _, p := range current {
				if p == newPath {
					found = true
					break
				}
			}

			if !found {
				current = append(current, newPath)
				s.manager.Publish(uapp.EventSettingsSaved, uapp.SettingsSavedPayload{
					Settings: state.Settings{
						CustomProgramPaths: current,
					},
				})

				s.manager.Publish(uapp.EventProgramsDiscoverRequested, nil)
				dialog.ShowInformation("Scan Path Added", newPath, w)
			}
		}, w)
	})

	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	content := container.NewBorder(
		container.NewHBox(refreshBtn, backBtn, addProgramBtn, addScanPathBtn),
		nil, nil, nil,
		s.programsList,
	)

	//Request initial discovery
	s.manager.Publish(uapp.EventProgramsDiscoverRequested, nil)

	return content
}
