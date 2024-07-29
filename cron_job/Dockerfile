FROM golang:1.22-alpine as build

WORKDIR /app

COPY ./go.mod go.mod
RUN go mod download
RUN go get github.com/PuerkitoBio/goquery

COPY ./cmd/parsers ./cmd/parsers
COPY ./cmd/filters ./cmd/filters
COPY ./cmd/types ./cmd/types
COPY ./cron_job/news_fetcher/ ./cron_job/news_fetcher
COPY ./cron_job/main.go main.go

RUN go build -o ./cron_job/job

FROM alpine:3.20

COPY --from=build ./app/cron_job/job ./news-fetching-job

ENTRYPOINT ["./news-fetching-job"]