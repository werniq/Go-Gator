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
ENV CERT_FILE=/app/cmd/server/certs/certificate.pem
ENV CERT_KEY=/app/cmd/server/certs/key.pem
ENV STORAGE_PATH=/tmp/

COPY --from=build /app/cmd/server/certs /cmd/server/certs
COPY --from=build /app/cmd/parsers/data/ $STORAGE_PATH
COPY --from=build /app/go-gator .
COPY --from=build $CERT_FILE $CERT_FILE
COPY --from=build $CERT_KEY $CERT_KEY

ENTRYPOINT /go-gator -p=$PORT -c=$CERT_FILE -k=$CERT_KEY -fs=$STORAGE_PATH