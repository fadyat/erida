ifneq (,$(wildcard ./.env))
	include .env
	export
endif

include helm/Makefile

lint:
	@golangci-lint run --issues-exit-code 1 \
		--print-issued-lines=true \
		  --config .golangci.yml ./...

run:
	@go run ./cmd/erida/main.go

stress:
	@go run ./cmd/stress/main.go

integration:
	@go test -v ./... -tags=integration


