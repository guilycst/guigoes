package domain

import (
	"github.com/golang-module/carbon/v2"
)

type Post struct {
	Dir       string
	Name      string
	Metadata  *Metadata
	Content   []byte
	Hash      string
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
}

type Metadata struct {
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
}
