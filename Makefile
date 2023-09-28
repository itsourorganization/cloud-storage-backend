.PHONY: run ling test docs gen

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go generate internal/repository/posgres/ent/
	go build -o app cmd/cloud/main.go && ./app

gen:
	go generate ./...
lint:
	golangci-lint run

test:
	go test ./...

docs:
	swag init --parseDependency -d cmd/cloud/,internal/transport/http-server/handlers/