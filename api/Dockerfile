# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./main.go

# Run stage
FROM alpine:latest

RUN adduser -D ymir
USER ymir

WORKDIR /home/ymir

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]