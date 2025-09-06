package app

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
)

func (a *App) WireProgramsHandlers() {
	discoverCh := a.Bus.Subscribe(EventProgramsDiscoverRequested)
	go func() {
		for range discoverCh {
			if a.ProgramService == nil {
				a.buildClients()
			}

			progs, err := a.ProgramService.DiscoverPrograms()
			if err != nil {
				continue

			}

			items := make([]ProgramItem, 0, len(progs))
			for _, p := range progs {
				items = append(items, ProgramItem{
					ID:          p.ID,
					Name:        p.Name,
					Path:        p.Path,
					Description: p.Description,
				})
			}

			a.Bus.Publish(EventProgramsDisocvered, ProgramDiscoveredPayload{
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

			hostID := a.State.Get().HostInfo.ID

			req := programs.CreateProgramRequest{
				Name:        payload.Program.Name,
				Path:        payload.Program.Path,
				Description: payload.Program.Description,
			}

			prog, err := a.ProgramService.RegisterProgram(context.Background(), req, hostID)
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
