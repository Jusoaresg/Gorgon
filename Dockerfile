# Stage 1: Build the backend
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix netgo -o /app/main .

# Stage 2: Create the final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

ENV IN_DOCKER=true

RUN mkdir -p /downloads /shows /config

EXPOSE 8080
CMD ["/app/main"]
