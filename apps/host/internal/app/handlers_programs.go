package app

import (
	"log"

	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
)

func (a *App) WireProgramsHandlers() {
	discoverCh := a.Bus.Subscribe(EventProgramsDiscoverRequested)
	go func() {
		for range discoverCh {
			if a.ProgramService == nil {
				a.buildClients()
			}
			if err := a.ProgramService.DiscoverAndSavePrograms(a.State.Get().Settings.CustomProgramPaths); err != nil {
				log.Printf("Failed to discover and save programs: %v", err)
			}

			list, err := a.ProgramService.GetLocalPrograms()
			if err != nil {
				continue
			}

			items := make([]ProgramItem, 0, len(list))
			for _, p := range list {
				items = append(items, ProgramItem{
					ID:          p.ID,
					Name:        p.Name,
					Path:        p.Path,
					Description: p.Description,
				})
			}

			a.Bus.Publish(EventProgramsDisocvered, ProgramsDiscoveredPayload{
				Programs: items,
			})
		}
	}()

	registerCh := a.Bus.Subscribe(EventProgramRegisterRequested)
	go func() {
		for evt := range registerCh {
			payload, ok := evt.(ProgramRegisterRequestedPayload)
			if !ok {
				continue
			}

			if a.ProgramService == nil {
				a.buildClients()
			}

			req := programs.CreateProgramRequest{
				Name:        payload.Program.Name,
				Path:        payload.Program.Path,
				Description: payload.Program.Description,
				HostID:      a.State.Get().HostInfo.ID,
			}

			prog, err := a.ProgramService.SaveProgram(req)
			if err != nil {
				continue
			}

			a.Bus.Publish(EventProgramRegistered, ProgramRegisteredPayload{
				Program: ProgramItem{
					ID:          prog.ID,
					Name:        prog.Name,
					Path:        prog.Path,
					Description: prog.Description,
				},
			})

			//refresh programs list
			a.Bus.Publish(EventProgramsDiscoverRequested, nil)
		}
	}()
}
