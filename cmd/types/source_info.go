package types

// Feed is a struct which is used to parse information about source: Name, Format and endpoint
// Name is basically name of the source
// Format is used to check what parsers should be used for that source
// Endpoint this field will be used to dynamically parse articles from that source
type Feed struct {
	Name     string `json:"name"`
	Format   string `json:"format"`
	Endpoint string `json:"endpoint"`
}
