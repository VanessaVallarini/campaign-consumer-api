run:
	go run ./cmd/ads-data-ingestor/main.go

clean:
	go mod tidy

config-local:
	./config.sh local

config-sandbox:
	./config.sh sandbox

config-production:
	./config.sh production

.PHONY: build run compose-up compose-down compose-infra-down compose-infra-up
compose-infra-up:
	docker-compose -f ./docker-compose.yml --profile infra up -d
compose-infra-down:
	docker-compose -f ./docker-compose.yml --profile infra down

air:
	air

air-init:
	air init

test:
	go test -count=1 ./... -covermode=count -coverprofile=coverage.out

mocks:
.PHONY: mockery
mockery:
	cd /tmp && go install github.com/vektra/mockery/v2@v2.42.0

mocks: mockery
	mockery --all --dir ./internal/service --output ./internal/service/mocks --with-expecter
	mockery --all --dir ./internal/client --output ./internal/client/mocks --with-expecter
	mockery --all --dir ./internal/listener --output ./internal/listener/mocks --with-expecter