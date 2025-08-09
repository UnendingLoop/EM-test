FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o em-test ./cmd/main.go
FROM alpine:latest
COPY --from=builder /app/em-test .
EXPOSE 8080
CMD ["./em-test"]
