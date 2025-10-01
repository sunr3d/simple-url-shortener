FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o simple-url-shortener ./cmd/main.go

FROM alpine:3.21

WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/simple-url-shortener .
RUN chown appuser:appuser /app
USER appuser

EXPOSE 8080
CMD ["./simple-url-shortener"]