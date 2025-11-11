.PHONY: build-dev mock go-doc

#TODO: build-production command that builds the app using the fyne-cli

go-doc:
	go doc -http=:6000

#builds a dev version with go-build instead of fyne package
build-dev:
	CGO_ENABLED=1 go build -o dist/host-app ./cmd/main.go

#generate mocks via go-mock and mockgen
mock:
	mockgen -package mocksession -destination internal/session/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/session Service
	mockgen -package mockhttpclient -destination pkg/httpclient/mocks/token_getter.go github.com/m1thrandir225/imperium/apps/host/pkg/httpclient TokenGetter
	mockgen -package mockhttpclient -destination pkg/httpclient/mocks/token_refresher.go github.com/m1thrandir225/imperium/apps/host/pkg/httpclient TokenRefresher
	mockgen -package mockprograms -destination internal/programs/mocks/db.go github.com/m1thrandir225/imperium/apps/host/internal/programs Database
	mockgen -package mockprograms -destination internal/programs/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/programs Service
	mockgen -package mockevent_broker -destination internal/events/mocks/broker.go github.com/m1thrandir225/imperium/apps/host/internal/events EventBroker
	mockgen -package mocktoken_refresher -destination internal/tokenrefresher/mocks/refresher.go github.com/m1thrandir225/imperium/apps/host/internal/tokenrefresher Refresher
	mockgen -package mockauth -destination internal/auth/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/auth Service
	mockgen -package mockhost -destination internal/host/mocks/status_manager.go github.com/m1thrandir225/imperium/apps/host/internal/host StatusManager
	mockgen -package mockui -destination internal/ui/mocks/screen.go github.com/m1thrandir225/imperium/apps/host/internal/ui Screen
	mockgen -package mockwebrtc -destination internal/webrtc/mocks/streamer.go github.com/m1thrandir225/imperium/apps/host/internal/webrtc Streamer
	mockgen -package mockrawg -destination pkg/rawg/mocks/client.go github.com/m1thrandir225/imperium/apps/host/pkg/rawg Client
