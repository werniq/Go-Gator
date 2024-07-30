# news_fetcher - Cron Job for fetching news

## Overview

The `news_fetcher` package provides functionality to fetch news articles from various sources. 
It is created to run a cron job that regularly fetches news articles, parses them, and stores them in JSON files.
The core component of this package is the `NewsFetchingJob` struct, 
which handles the main job of fetching and parsing news.

## Implemented Features

1. **NewsFetchingJob**: Represents a job for fetching news articles with a timestamp.
2. **RunJob Function**: Initializes and runs a `NewsFetchingJob` struct
3. **Execute Method**: Fetches news, parses it, and writes the parsed data to a JSON file named with the
current date in the format `YYYY-MM-DD`.

## Usage

### Using Golang

   ```sh
   go build -o ./bin/news_fetcher
   ```

### Using docker

```sh
docker build -t {IMAGE_TAG_TITLE} .
```