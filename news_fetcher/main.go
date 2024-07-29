package main

import (
	"gogator/news_fetcher/fetchnews"
	"log"
)

func main() {
	err := fetchnews.RunJob()
	if err != nil {
		log.Fatalln(err)
	}
}
