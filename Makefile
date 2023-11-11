ifneq (,$(wildcard ./.env))
	include .env
	export
endif

include helm/Makefile

run:
	@go run ./cmd/erida/main.go

stress:
	@go run ./cmd/stress/main.go

integration:
	@go test -v ./... -tags=integration


