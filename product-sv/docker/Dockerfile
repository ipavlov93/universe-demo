FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/product-sv

RUN go build -o product-sv ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/product-sv/product-sv .

ENTRYPOINT ["/app/product-sv"]
