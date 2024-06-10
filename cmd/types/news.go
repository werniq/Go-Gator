package types

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Json struct {
	Articles []News `json:"articles"`
}

type Channel struct {
	Items []News `xml:"item"`
}

type News struct {
	Title       string `json:"title" xml:"title"`
	PubDate     string `json:"publishedAt" xml:"pubDate"`
	Description string `json:"description" xml:"description"`
	Publisher   string `xml:"source"`
	Link        string `json:"url" xml:"link"`
}
