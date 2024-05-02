package services

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/pkg"
)

type LocalPostService struct{}

func NewLocalPostService() ports.PostService {
	return &LocalPostService{}
}

var defaultOpts = &ports.PostsOptions{
	FetchN:  pkg.Ptr(uint64(10)),
	SkipN:   pkg.Ptr(uint64(0)),
	Content: pkg.Ptr(ports.None),
}

func (lps LocalPostService) Posts(opts *ports.PostsOptions) ([]*domain.Post, error) {
	opts = resolvePostsOpts(opts)

	metaPaths, err := filepath.Glob(pkg.POSTS_PATH + "**/metadata.json")
	if err != nil {
		return nil, err
	}

	metas := []*domain.Metadata{}
	for _, metaPath := range metaPaths {
		meta, err := getMeta(metaPath)
		if err != nil {
			return nil, err
		}

		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].CreatedAt.Gt(metas[j].CreatedAt.Carbon)
	})

	skip := *opts.SkipN + *opts.FetchN
	if skip > uint64(len(metas)) {
		skip = uint64(len(metas))
	}
	metas = metas[*opts.SkipN:skip]
	posts := []*domain.Post{}
	for _, meta := range metas {
		post, err := getPost(meta, opts.Content)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (lps LocalPostService) GetPost(postName string) (*domain.Post, error) {
	return getPostWithContent(postName, pkg.Ptr(ports.HTML))
}

func getPostWithContent(postName string, opt *ports.PostsContentOpt) (*domain.Post, error) {
	meta, err := getMeta(pkg.POSTS_PATH + postName + "/metadata.json")
	if err != nil {
		return nil, err
	}

	post, err := getPost(meta, opt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (lps LocalPostService) GetPostAsset(postName string, assetName string) (string, error) {
	var postAssetPath = pkg.POSTS_PATH + postName + "/assets/" + assetName
	slog.Info("Serving asset", postAssetPath, "")
	if err := checkFileExists(postAssetPath); err != nil {
		return "", err
	}
	return postAssetPath, nil
}

func (lps LocalPostService) SearchPosts(term string) ([]*domain.Post, error) {
	index, err := bleve.Open(pkg.BLEVE_IDX_PATH)
	if err != nil {
		return nil, err
	}
	defer index.Close()

	query := bleve.NewMatchQuery(term)
	query.Fuzziness = 2
	search := bleve.NewSearchRequest(query)
	result, err := index.Search(search)
	if err != nil {
		return nil, err
	}

	if result.Total == 0 {
		return []*domain.Post{}, nil
	}

	hits := result.Hits[:10]
	posts := []*domain.Post{}
	for _, hit := range hits {
		if hit == nil {
			continue
		}
		post, err := getPostWithContent(path.Base(hit.ID), pkg.Ptr(ports.None))
		if err != nil {
			log.Printf("Post %s not found, reindex migh be necessary\n", hit.ID)
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func getPostContent(contentPath string, contentOpt *ports.PostsContentOpt) ([]byte, error) {
	if contentOpt == nil || *contentOpt == ports.None {
		return []byte{}, nil
	}

	content, err := os.ReadFile(contentPath)
	if err != nil {
		return nil, err
	}

	if *contentOpt == ports.Markdown {
		return content, nil
	}

	return markdownToHTML(content), nil
}

func resolvePostsOpts(opts *ports.PostsOptions) *ports.PostsOptions {
	if opts == nil {
		return defaultOpts
	} else {
		if opts.FetchN == nil {
			opts.FetchN = defaultOpts.FetchN
		}
		if opts.SkipN == nil {
			opts.SkipN = defaultOpts.SkipN
		}
		if opts.Content == nil {
			opts.Content = defaultOpts.Content
		}
	}
	return opts
}

func setPostGitTrackingInfo(postName string, meta *domain.Metadata) error {
	var gitTrack = pkg.POSTS_PATH + postName + "/git-log.track"
	track, err := os.ReadFile(gitTrack)
	if err != nil || len(track) == 0 {
		meta.CreatedAt = carbon.Now().ToDateTimeStruct()
		meta.UpdatedAt = carbon.Now().ToDateTimeStruct()
		slog.Error("Error reading git-log.track file: ", err)
		return nil
	}

	lines := strings.Split(string(track), "\n")
	first := lines[0]
	firstCmp := strings.Split(first, "_")
	ft := carbon.Parse(firstCmp[0]).ToDateTimeStruct()

	//last line is empty
	if len(lines) == 2 {
		meta.CreatedAt = ft
		meta.UpdatedAt = ft
		meta.Hash = firstCmp[1]
		return nil
	}

	last := lines[len(lines)-2]
	lastCmp := strings.Split(last, "_")
	lt := carbon.Parse(lastCmp[0]).ToDateTimeStruct()

	meta.CreatedAt = lt
	meta.UpdatedAt = ft
	meta.Hash = lastCmp[1]

	return nil
}

func getPost(meta *domain.Metadata, contentOpt *ports.PostsContentOpt) (*domain.Post, error) {
	dir := pkg.POSTS_PATH + meta.PostName + "/body.md"
	content, err := getPostContent(dir, contentOpt)
	if err != nil {
		return nil, err
	}

	post := &domain.Post{
		Dir:      "/posts/" + meta.PostName,
		Name:     meta.PostName,
		Metadata: meta,
		Content:  content,
	}
	return post, nil
}

func checkFileExists(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return &domain.FSResourceNotFoundError{Msg: filePath, Err: err}
	}
	return err
}

func getMeta(metaPath string) (*domain.Metadata, error) {
	if err := checkFileExists(metaPath); err != nil {
		return nil, err
	}

	postName := filepath.Base(filepath.Dir(metaPath))
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	meta := &domain.Metadata{
		PostName: postName,
	}
	err = json.Unmarshal(metaBytes, meta)
	if err != nil {
		slog.Error("Invalid metadata.json: ", metaPath, err)
	}

	err = setPostGitTrackingInfo(postName, meta)
	if err != nil {
		return nil, err
	}
	meta.Thumb = path.Join("/assets", "thumb_"+path.Base(meta.Cover))
	return meta, nil
}
