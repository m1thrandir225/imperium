package app

import (
	"log"

	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
)

func (a *App) WireSettingsHandlers() {
	ch := a.Bus.Subscribe(EventSettingsSaved)
	go func() {
		for evt := range ch {
			payload, ok := evt.(SettingsSavedPayload)
			if !ok {
				continue
			}

			err := a.State.Update(func(s *state.AppState) {
				if payload.Settings.ServerAddress != "" {
					s.Settings.ServerAddress = payload.Settings.ServerAddress
				}
				if payload.Settings.FFmpegPath != "" {
					s.Settings.FFmpegPath = payload.Settings.FFmpegPath
				}

				if payload.Settings.Encoder != "" {
					s.Settings.Encoder = payload.Settings.Encoder
				}
				if payload.Settings.Framerate != 0 {
					s.Settings.Framerate = payload.Settings.Framerate
				}

				if payload.Settings.Bitrate != "" {
					s.Settings.Bitrate = payload.Settings.Bitrate
				}

				if len(payload.Settings.CustomProgramPaths) > 0 {
					s.Settings.CustomProgramPaths = payload.Settings.CustomProgramPaths
				}
			})
			if err != nil {
				log.Printf("failed to update state: %v", err)
			}

			a.AuthBaseURL = a.State.Get().Settings.ServerAddress

			a.buildClients()
			if a.SessionService != nil {
				st := a.State.Get()
				a.SessionService.UpdateVideoConfig(&video.Config{
					Encoder:    st.Settings.Encoder,
					FPS:        st.Settings.Framerate,
					FFMPEGPath: st.Settings.FFmpegPath,
				})
			}
			a.Bus.Publish(EventStateSaved, a.State.Get())
		}
	}()
}
