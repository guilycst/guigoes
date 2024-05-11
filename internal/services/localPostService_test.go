package services_test

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"

	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
	_ "github.com/guilycst/guigoes/testing"
)

func TestPostShouldFetchN(t *testing.T) {
	pkg.POSTS_PATH = "./posts/"
	opts := ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1)),
		SkipN:   pkg.Ptr(uint64(0)),
		Content: pkg.Ptr(ports.None),
	}

	lps := services.NewLocalPostService()
	posts, err := lps.Posts(&opts)
	if err != nil {
		t.Fatalf("Error fetching posts: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("Expected 1 post, got %d", len(posts))
	}
}

func TestPostShouldFetchNoneWrongPostPath(t *testing.T) {
	pkg.POSTS_PATH = "./wrong_posts_dir/"
	opts := ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1)),
		SkipN:   pkg.Ptr(uint64(0)),
		Content: pkg.Ptr(ports.None),
	}

	lps := services.NewLocalPostService()
	posts, err := lps.Posts(&opts)
	if err != nil {
		t.Fatalf("Error fetching posts: %v", err)
	}

	if len(posts) != 0 {
		t.Fatalf("Expected 0 post, got %d", len(posts))
	}
}

func TestPostShouldFetchNoContentOnlyMetadata(t *testing.T) {
	pkg.POSTS_PATH = "./posts/"
	opts := ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1)),
		SkipN:   pkg.Ptr(uint64(0)),
		Content: pkg.Ptr(ports.None),
	}

	lps := services.NewLocalPostService()
	posts, err := lps.Posts(&opts)
	if err != nil {
		t.Fatalf("Error fetching posts: %v", err)
	}

	post := posts[0]
	if len(post.Content) != 0 {
		t.Fatalf("Expected empty content, got %s", post.Content)
	}
}

func TestPostShouldFetchRenderedHTML(t *testing.T) {
	pkg.POSTS_PATH = "./posts/"
	opts := ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1)),
		SkipN:   pkg.Ptr(uint64(0)),
		Content: pkg.Ptr(ports.HTML),
	}

	lps := services.NewLocalPostService()
	posts, err := lps.Posts(&opts)
	if err != nil {
		t.Fatalf("Error fetching posts: %v", err)
	}

	post := posts[0]
	if len(post.Content) == 0 {
		t.Fatalf("Expected non-empty content, got empty")
	}

	err = tryParseHTML(post.Content)
	if err != nil {
		t.Fatalf("Error parsing HTML: %v", err)
	}
}

func tryParseHTML(data []byte) error {
	r := bytes.NewReader(data)
	d := xml.NewDecoder(r)

	// Configure the decoder for HTML; leave off strict and autoclose for XHTML
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity
	for {
		_, err := d.Token()
		switch err {
		case io.EOF:
			return nil // We're done, it's valid!
		case nil:
		default:
			return err // Oops, something wasn't right
		}
	}
}
