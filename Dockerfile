# Stage 1: Build stage
FROM golang:1.21-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o go-gator .

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

ENV APP_MODE=DOCKER

RUN echo APP_MODE=DOCKER >> .env
RUN mkdir -p ./cmd/parsers/data/ && mkdir -p ./cmd/server/certs/

# Copy the binary from the build stage
COPY --from=build /app/go-gator .
COPY --from=build /app/cmd/server/certs/certificate.pem ./cmd/server/certs/certificate.pem
COPY --from=build /app/cmd/server/certs/key.pem ./cmd/server/certs/key.pem

# Set the timezone and install CA certificates
# RUN apk --no-cache add ca-certificates tzdata

EXPOSE 8080

# Set the entrypoint command
ENTRYPOINT ["/app/go-gator"]