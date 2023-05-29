FROM golang:1.20-alpine


RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go install github.com/cosmtrek/air@latest
COPY go.mod go.sum ./

RUN apk update && apk add --no-cache build-base

RUN go mod download
RUN go build -o booking ./cmd/web/*

CMD ["air", "-c", ".air.toml", "./booking"]