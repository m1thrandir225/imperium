.PHONY: mock

mock:
	mockgen -package mocksession -destination internal/session/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/session Service
	mockgen -package mockhttpclient -destination internal/httpclient/mocks/token_getter.go github.com/m1thrandir225/imperium/apps/host/internal/httpclient TokenGetter
	mockgen -package mockhttpclient -destination internal/httpclient/mocks/token_refresher.go github.com/m1thrandir225/imperium/apps/host/internal/httpclient TokenRefresher
	#mockgen -package mockprograms -destination internal/programs/mocks/db.go github.com/m1thrandir225/imperium/apps/host/internal/programs Database
	#mockgen -package mockprograms -destination internal/programs/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/programs Service
	mockgen -package mockevent_broker -destination internal/event_broker/mocks/broker.go github.com/m1thrandir225/imperium/apps/host/internal/event_broker EventBroker
	mockgen -package mocktoken_refresher -destination internal/token_refresher/mocks/refresher.go github.com/m1thrandir225/imperium/apps/host/internal/token_refresher Refresher
	mockgen -package mockauth -destination internal/auth/mocks/service.go github.com/m1thrandir225/imperium/apps/host/internal/auth Service
