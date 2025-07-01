# Stage 1: Build the frontend
FROM node:20-alpine AS frontend
WORKDIR /app/assets/front
COPY assets/front/package.json assets/front/package-lock.json ./
RUN npm ci
COPY assets/front/. .
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/assets/front/build ./assets/front/build
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

ENV IN_DOCKER=true

RUN mkdir -p /downloads /shows /config

EXPOSE 8080
CMD ["/app/main"]

