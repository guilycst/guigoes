package ports

import (
	"github.com/guilycst/guigoes/internal/core/domain"
)

type PostService interface {
	Posts(opts *PostsOptions) ([]*domain.Post, error)
	GetPost(postName string) (*domain.Post, error)
	GetPostAsset(postName string, assetName string) (string, error)
	SearchPosts(term string) ([]*domain.Post, error)
}

type PostsContentOpt string

// I really miss real enums
func (c PostsContentOpt) String() string {
	return string(c)
}

const (
	HTML     PostsContentOpt = "html"
	Markdown PostsContentOpt = "markdown"
	None     PostsContentOpt = "none"
)

type PostsOptions struct {
	FetchN  *uint64
	SkipN   *uint64
	Content *PostsContentOpt
}
