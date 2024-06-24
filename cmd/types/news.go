package types

// RSS struct is used to parse articles in RSS format.
// Because each resource has its own data output format,
// this model will be used when we have the following structure:
//
// <rss version="2.0">
//
//	 <channel>
//		  News fields...
//	 </channel>
//
// </rss>
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Json struct {
	Articles []News `json:"articles"`
}

type Channel struct {
	Items []News `xml:"item"`
}

// News is one of the main models in news aggregator.
// It has few fields inside:
// /   1. Title			- Headline of the article
// /   2. Description 	- Description of the article
// /   3. PubDate 		- Publication Date
// /   4. Link 			- Link to the article
// /   5. Publisher 	- Optional: Author or publisher of the article
//
// It will be used through the application for different operations, such as:
//  1. Parsing
//  2. Logging
type News struct {
	Title       string `json:"title" xml:"title"`
	PubDate     string `json:"publishedAt" xml:"pubDate"`
	Description string `json:"description" xml:"description"`
	Publisher   string `xml:"source" json:"Publisher"`
	Link        string `json:"url" xml:"link"`
}
