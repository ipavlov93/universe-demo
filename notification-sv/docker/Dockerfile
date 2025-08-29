FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/notification-sv

RUN go build -o notification-sv ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/notification-sv/notification-sv .

ENTRYPOINT ["/app/notification-sv"]
