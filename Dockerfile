# Dockerfile
# This file builds the Go API application into a runnable container image.
# It is used to package the app so it can run in Docker independently of the host machine.
# The image is then used by Docker Compose or by a docker run command.

# stage: builder
FROM golang:1.25.5-alpine AS builder
WORKDIR /src

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static, stripped binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /app/ecommerce ./cmd

# runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app

# copy binary from builder
COPY --from=builder /app/ecommerce /app/ecommerce

# create non-root user and drop privileges
RUN addgroup -S app && adduser -S app -G app
USER app

# runtime env & port
EXPOSE 8080

ENTRYPOINT ["./ecommerce"]