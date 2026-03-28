FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/fasthttp

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/server .
COPY config/config.docker.yaml ./config/config.docker.yaml

ENV CONFIG_PATH=/app/config/config.docker.yaml

EXPOSE 8080

CMD ["./server"]
