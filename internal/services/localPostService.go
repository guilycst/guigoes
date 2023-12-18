package services

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-module/carbon/v2"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/pkg"
)

type LocalPostService struct{}

func NewLocalPostService() ports.PostService {
	return &LocalPostService{}
}

func (lps LocalPostService) Index() ([]*domain.Post, error) {
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
		posts[dir] = &domain.Post{Dir: dir, Name: filepath.Base(filepath.Dir(md))}
	}

	var validPosts = []*domain.Post{}
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

		err = setPostGitTrackingInfo(post)
		if err != nil {
			return nil, err
		}
		validPosts = append(validPosts, post)
	}

	sort.Slice(validPosts, func(i, j int) bool {
		return validPosts[i].CreatedAt.Gt(validPosts[j].CreatedAt.Carbon)
	})

	return validPosts, nil
}

func (lps LocalPostService) GetPost(postName string) (*domain.Post, error) {
	var postMd = pkg.POSTS_PATH + postName + "/body.md"
	var postMeta = pkg.POSTS_PATH + postName + "/metadata.json"
	var post = &domain.Post{
		Dir:  filepath.Dir(postMd),
		Name: postName,
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

	err = setPostGitTrackingInfo(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func setPostGitTrackingInfo(post *domain.Post) error {
	var gitTrack = pkg.POSTS_PATH + post.Name + "/git-log.track"
	track, err := os.ReadFile(gitTrack)
	if err != nil {
		return err
	}

	lines := strings.Split(string(track), "\n")
	first := lines[0]
	firstCmp := strings.Split(first, "_")
	ft := carbon.Parse(firstCmp[0]).ToDateTimeStruct()

	//last line is empty
	if len(lines) == 2 {
		post.CreatedAt = ft
		post.UpdatedAt = ft
		post.Hash = firstCmp[1]
		return nil
	}

	last := lines[len(lines)-2]
	lastCmp := strings.Split(last, "_")
	lt := carbon.Parse(lastCmp[0]).ToDateTimeStruct()

	post.CreatedAt = lt
	post.UpdatedAt = ft
	post.Hash = lastCmp[1]

	return nil
}

func (lps LocalPostService) GetPostAsset(postName string, assetName string) (string, error) {
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
