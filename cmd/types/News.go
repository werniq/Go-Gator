package types

import (
	"encoding/xml"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Image       Image  `xml:"image"`
	Items       []News `xml:"item"`
}

type Image struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Link  string `xml:"link"`
}

type News struct {
	Title       string      `json:"title" xml:"title"`
	Link        string      `json:"link" xml:"link"`
	Guid        string      `json:"guid" xml:"guid"`
	LinkedVideo string      `json:"linkedVideo" xml:"LinkedVideo"`
	PubDate     string      `json:"pubDate" xml:"pubDate"`
	Description string      `json:"description" xml:"description"`
	Category    string      `json:"category" xml:"category"`
	Thumbnails  []Thumbnail `json:"thumbnails" xml:"thumbnail"`
}

type Thumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}
