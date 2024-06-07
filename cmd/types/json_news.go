package types

import "strings"

type Json struct {
	Status       string     `json:"status"`
	TotalResults int        `json:"totalResults"`
	Articles     []JsonNews `json:"articles"`
}

type JsonNews struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	UrlToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func JsonNewsToNews(jsonNews []JsonNews, source string) []News {
	news := []News{}
	for _, article := range jsonNews {
		news = append(news, News{
			Title:       article.Title,
			Link:        article.Url,
			PubDate:     article.PublishedAt,
			Description: article.Description,
			Source:      strings.Split(source, ".")[0],
		})
	}

	return news
}
