FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/filters ./cmd/filters
COPY ./cmd/parsers ./cmd/parsers
COPY ./cmd/parsers/data ./cmd/parsers/data
COPY ./cmd/templates ./cmd/templates
COPY ./cmd/types ./cmd/types
COPY ./cmd/validator ./cmd/validator
COPY ./cmd/server ./cmd/server
COPY ./main.go ./main.go

RUN go build -o go-gator .

FROM alpine:3.20

ENV PORT=443
ENV STORAGE_PATH=./data

COPY --from=build /app/cmd/server/certs ./cmd/server/certs
COPY --from=build /app/cmd/parsers/data $STORAGE_PATH
COPY --from=build /app/go-gator .

ENTRYPOINT /go-gator -p=$PORT -fs=$STORAGE_PATH