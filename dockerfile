# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o crawler ./webcrawler/main.go

RUN chmod +x crawler

# Run stage
FROM alpine:latest

WORKDIR /runtime

COPY --from=builder /app/crawler ./crawler

CMD ["./crawler"]