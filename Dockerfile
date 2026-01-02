FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shortlink ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/shortlink ./shortlink
COPY configs/config.example.yaml /app/configs/config.yaml
EXPOSE 8080
ENV APP_CONFIG=/app/configs/config.yaml
ENTRYPOINT ["./shortlink"]
