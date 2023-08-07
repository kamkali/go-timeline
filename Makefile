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
	docker build -t github.com/kamkali/AiTSI_projekt .
#	docker tag go-timeline:latest github.com/kamkali/AiTSI_projekt && \

push_image:
	heroku container:login && \
 	heroku container:push -a timeline-backend web

deploy:
	heroku container:release -a timeline-backend web

logs:
	heroku logs -a timeline-backend --tail

.PHONY: compose tools generate lint build_image push_image deploy logs