# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/worker .

CMD ["./api"]