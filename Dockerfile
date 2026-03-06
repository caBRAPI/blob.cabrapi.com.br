# 1 - Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o blob main.go

# 2 - Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/blob /app/blob
EXPOSE 3000
CMD ["/app/blob"]
