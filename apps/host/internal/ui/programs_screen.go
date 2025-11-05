package ui

import (
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type ProgramsScreen struct {
	manager      *uiManager
	programsList *widget.List
	programs     []uapp.ProgramItem
	subscribed   bool
}

func NewProgramsScreen(manager *uiManager) *ProgramsScreen {
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
				payload, ok := evt.(uapp.ProgramsDiscoveredPayload)
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
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			program := s.programs[id]

			nameLabel := box.Objects[0].(*widget.Label)
			pathLabel := box.Objects[1].(*widget.Label)

			nameLabel.SetText(program.Name)
			pathLabel.SetText(util.ShortPath(program.Path))
		},
	)

	refreshBtn := widget.NewButton("Refresh Programs", func() {
		s.manager.publish(uapp.EventProgramsDiscoverRequested, nil)
	})

	addProgramBtn := widget.NewButton("Add Program", func() {
		fd := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
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

			s.manager.publish(uapp.EventProgramRegisterRequested, uapp.ProgramRegisterRequestedPayload{
				Program: uapp.ProgramItem{
					Name:        name,
					Path:        path,
					Description: "",
				},
			})
		}, w)
		home, _ := os.UserHomeDir()
		u := storage.NewFileURI(home)
		l, _ := storage.ListerForURI(u)
		fd.SetLocation(l)
		fd.Show()
	})

	addScanPathBtn := widget.NewButton("Add Scan Path", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}

			current := s.manager.GetState().Settings.CustomProgramPaths
			newPath := uri.Path()

			found := slices.Contains(current, newPath)

			if !found {
				current = append(current, newPath)
				s.manager.publish(uapp.EventSettingsSaved, uapp.SettingsSavedPayload{
					Settings: state.Settings{
						CustomProgramPaths: current,
					},
				})

				s.manager.publish(uapp.EventProgramsDiscoverRequested, nil)
				dialog.ShowInformation("Scan Path Added", newPath, w)
			}
		}, w)

		home, _ := os.UserHomeDir()
		u := storage.NewFileURI(home)
		l, _ := storage.ListerForURI(u)
		fd.SetLocation(l)
		fd.Show()
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
	s.manager.publish(uapp.EventProgramsDiscoverRequested, nil)

	return content
}
