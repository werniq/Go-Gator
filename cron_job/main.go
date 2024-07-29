package main

import (
	"gogator/cron_job/news_fetcher"
	"log"
)

func main() {
	err := news_fetcher.RunJob()
	if err != nil {
		log.Fatalln(err)
	}
}
