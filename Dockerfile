FROM golang:1.19-alpine as build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN cd ./cmd && CGO_ENABLED=0 go build -o app

FROM alpine as runner

COPY --from=build /app/.env /
COPY --from=build /app/cmd/app /

EXPOSE 8080

ENTRYPOINT ["/app"]