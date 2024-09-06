FROM golang:1.22-alpine as build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd/parsers ./cmd/parsers
COPY ./cmd/filters ./cmd/filters
COPY ./cmd/types ./cmd/types
COPY ./news_fetcher/ ./news_fetcher
COPY ./news_fetcher/main.go main.go
COPY ./news_fetcher/fetch_news_job.go fetch_news_job.go

RUN go build -o ./news_fetcher_job .

FROM alpine:3.20

ENV FILE_STORAGE=/tmp/

COPY --from=build ./app/news_fetcher_job ./news-fetching-job

ENTRYPOINT ./news-fetching-job -fs=${FILE_STORAGE}