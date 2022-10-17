compose:
	docker compose up -d

tools: install-mockery

install-mockery:
	go install github.com/vektra/mockery/v2@latest

generate:
	go generate ./...

lint:
	golangci-lint run ./...


.PHONY: compose tools generate lint