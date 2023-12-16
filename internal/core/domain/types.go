package domain

type Post struct {
	Dir      string
	Metadata *Metadata
}

type Metadata struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
