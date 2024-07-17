FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-gator .

FROM alpine:3.20

WORKDIR /app

RUN mkdir -p ./cmd/parsers/data/ && mkdir -p ./cmd/server/certs/

COPY --from=build /app/go-gator .
COPY --from=build /app/cmd/server/certs/certificate.pem ./cmd/server/certs/certificate.pem
COPY --from=build /app/cmd/server/certs/key.pem ./cmd/server/certs/key.pem
COPY --from=build /app/cmd/parsers/data ./cmd/parsers/data

ENTRYPOINT ["/app/go-gator"]