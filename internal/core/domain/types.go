package domain

type Post struct {
	Dir      string
	Metadata *Metadata
	Content  []byte
}

type Metadata struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
