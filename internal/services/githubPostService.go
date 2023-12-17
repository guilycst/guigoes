package services

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/pkg"
)

type GithubPostService struct{}

func NewGithubPostService() ports.PostService {
	return &LocalPostService{}
}

func (lps GithubPostService) Index() (map[string]*domain.Post, error) {
	mds, err := filepath.Glob(pkg.POSTS_PATH + "**/*.md")
	if err != nil {
		return nil, err
	}

	metas, err := filepath.Glob(pkg.POSTS_PATH + "**/metadata.json")
	if err != nil {
		return nil, err
	}

	posts := make(map[string]*domain.Post)
	for _, md := range mds {
		dir := "/posts/" + filepath.Base(filepath.Dir(md))
		posts[dir] = &domain.Post{Dir: dir}
	}

	for _, meta := range metas {
		dir := "/posts/" + filepath.Base(filepath.Dir(meta))
		post, ok := posts[dir]
		if !ok {
			log.Println("Dangling metadata.json: ", meta)
			continue
		}

		metaBytes, err := os.ReadFile(meta)
		if err != nil {
			return nil, err
		}

		post.Metadata = &domain.Metadata{}
		err = json.Unmarshal(metaBytes, post.Metadata)
		if err != nil {
			log.Println("Invalid metadata.json: ", meta, err)
		}
	}

	return posts, nil
}

func (lps GithubPostService) GetPost(postName string) (*domain.Post, error) {
	var postMd = pkg.POSTS_PATH + postName + "/body.md"
	var postMeta = pkg.POSTS_PATH + postName + "/metadata.json"
	var post = &domain.Post{
		Dir: filepath.Dir(postMd),
	}

	metaBytes, err := os.ReadFile(postMeta)
	if err != nil {
		return nil, err
	}

	post.Metadata = &domain.Metadata{}
	err = json.Unmarshal(metaBytes, post.Metadata)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(postMd)
	if err != nil {
		return nil, err
	}

	post.Content = markdownToHTML(content)
	return post, nil
}

func (lps GithubPostService) GetPostAsset(postName string, assetName string) (string, error) {
	var postAssetPath = pkg.POSTS_PATH + postName + "/assets/" + assetName
	log.Println("Serving asset: ", postAssetPath)
	if _, err := os.Stat(postAssetPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", &domain.AssetNotFoundError{}
		}
		return "", err
	}
	return postAssetPath, nil
}
