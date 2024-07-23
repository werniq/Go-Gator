FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/filters ./cmd/filters
COPY ./cmd/parsers ./cmd/parsers
COPY ./cmd/templates ./cmd/templates
COPY ./cmd/types ./cmd/types
COPY ./cmd/validator ./cmd/validator
COPY ./cmd/server ./cmd/server
COPY ./main.go ./main.go

RUN go build -o go-gator .

FROM alpine:3.20

WORKDIR /app

RUN mkdir -p ./cmd/parsers/data/ && mkdir -p ./cmd/server/certs/

ENV PORT=443
ENV UPDATES_FREQUENCY=4
ENV CERT_FILE="/app/cmd/server/certs/certificate.pem"
ENV CERT_KEY="/app/cmd/server/certs/key.pem"
ENV STORAGE_PATH="/app/cmd/parsers/data"


COPY --from=build /app/go-gator .
COPY --from=build /app/cmd/server/certs/certificate.pem ./cmd/server/certs/certificate.pem
COPY --from=build /app/cmd/server/certs/key.pem ./cmd/server/certs/key.pem
COPY --from=build /app/cmd/parsers/data ./cmd/parsers/data

ENTRYPOINT ["/app/go-gator", \
            "-p", "$PORT", \
            "-f", "$UPDATES_FREQUENCY", \
            "-c", "$CERT_FILE", \
            "-k", "$CERT_KEY", \
            "-fs", "$STORAGE_PATH"]
