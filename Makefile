compose:
	docker compose up -d

tools: install-mockery

install-mockery:
	go install github.com/vektra/mockery/v2@latest

generate:
	go generate ./...

test:
	go test -short -race ./...


lint:
	golangci-lint run ./...

#### Deployment
build_image:
	docker build -t github.com/kamkali/go-timeline .
#	docker tag go-timeline:latest github.com/kamkali/go-timeline && \

push_image:
	heroku container:login && \
 	heroku container:push -a go-timeline web

deploy:
	heroku container:release -a go-timeline web

logs:
	heroku logs -a go-timeline --tail

.PHONY: compose tools generate lint build_image push_image deploy logs