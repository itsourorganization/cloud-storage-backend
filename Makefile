.PHONY: run ling test docs

run:
	go build -o app cmd/cloud/main.go && CONFIG_PATH=config/config.yaml ./app

lint:
	golangci-lint run

test:
	go test ./...

docs:
	swag init --parseDependency -d cmd/cloud/,internal/transport/http-server/handlers/