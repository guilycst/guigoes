package domain

import (
	"github.com/golang-module/carbon/v2"
)

type Post struct {
	Dir      string
	Name     string
	Metadata *Metadata
	Content  []byte
}

type Metadata struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CreatedAt   carbon.DateTime
	UpdatedAt   carbon.DateTime
	Hash        string
	PostName    string
}
