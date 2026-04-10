FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY proto ./proto
COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env /build/.env
COPY deployment /build/deployment

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /build/app .

# Копируем .env и deployment из builder
COPY --from=builder /build/.env .
COPY --from=builder /build/deployment ./deployment

RUN chown -R appuser:appuser /app && chmod +x /app/app

USER appuser
EXPOSE 50051
CMD ["/app/app"]