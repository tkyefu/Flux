# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /flux

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /flux /app/flux
EXPOSE 8080
CMD ["/app/flux"]
