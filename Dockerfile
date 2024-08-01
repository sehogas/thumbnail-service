# ./Dockerfile

FROM golang:1.22.2-alpine AS builder
LABEL MAINTAINER "Sebastian Hogas <sehogas@gmail.com>"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

ARG VERSION=1.0.0
ENV FLAGS="-s -w -X 'main.Version=${VERSION}'"
RUN go build -ldflags="${FLAGS}" -o ./thumbnail-service ./cmd/api/main.go

FROM scratch

COPY --from=builder ["/build/thumbnail-service", "/"]

EXPOSE 3010

ENTRYPOINT ["/thumbnail-service"]