FROM golang:1.22-alpine as build

WORKDIR /app

COPY ./go.mod go.mod
RUN go mod download
RUN go get github.com/PuerkitoBio/goquery

COPY ./cmd/parsers ./cmd/parsers
COPY ./cmd/filters ./cmd/filters
COPY ./cmd/types ./cmd/types
COPY ./news_fetcher/fetchnews/ ./news_fetcher/fetchnews
COPY ./news_fetcher/main.go main.go

RUN go build -o ./news-fetching-job .

FROM alpine:3.20

COPY --from=build ./app/news-fetching-job ./news-fetching-job

ENTRYPOINT ["./news-fetching-job"]