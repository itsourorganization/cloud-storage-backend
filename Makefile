.PHONY: run ling

run:
	go build -o app cmd/cloud/main.go && CONFIG_PATH=config/config.yaml ./app

lint:
	golangci-lint run
