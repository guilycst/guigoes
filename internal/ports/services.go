package ports

import "github.com/guilycst/guigoes/internal/core/domain"

type PostService interface {
	Index() ([]*domain.Post, error)
	GetPost(postName string) (*domain.Post, error)
	GetPostAsset(postName string, assetName string) (string, error)
}
