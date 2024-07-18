FROM golang:1.21-alpine AS build

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
# Questions:
# Should taskfile be included?
# COPY Taskfile.yml Taskfile.yml
# Also, why we can't use `COPY . .` with dockerignore?

RUN go build -o go-gator .

FROM alpine:3.20

WORKDIR /app

RUN mkdir -p ./cmd/parsers/data/ && mkdir -p ./cmd/server/certs/

COPY --from=build /app/go-gator .
COPY --from=build /app/cmd/server/certs/certificate.pem ./cmd/server/certs/certificate.pem
COPY --from=build /app/cmd/server/certs/key.pem ./cmd/server/certs/key.pem
COPY --from=build /app/cmd/parsers/data ./cmd/parsers/data

ENTRYPOINT ["/app/go-gator"]