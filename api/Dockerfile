# syntax=docker/dockerfile:1

FROM golang:1.24.2-alpine

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.* .
RUN go mod download

COPY . .

# Build
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /ymir-api

EXPOSE 8080

# Run
CMD ["/ymir-api"]