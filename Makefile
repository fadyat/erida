ifneq (,$(wildcard ./.env))
	include .env
	export
endif

run:
	@go run ./cmd/main.go

integration:
	@go test -v ./... -tags=integration


