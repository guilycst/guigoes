package handlers

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/pkg"
	"github.com/guilycst/guigoes/web/templates"
)

func NewGin() *gin.Engine {
	r := gin.Default()
	r.GET("/", Index)
	r.GET("/posts/:post", Post)
	r.GET("/posts/assets/:asset", PostAsset)
	return r
}

func PostAsset(c *gin.Context) {
	ref := c.Request.Header.Get("Referer")
	if ref == "" {
		c.AbortWithStatus(404)
	}

	url, err := url.Parse(ref)
	if err != nil {
		c.AbortWithError(500, err)
	}

	assetName := c.Param("asset")
	var postAssetPath = pkg.POSTS_PATH + filepath.Base(url.Path) + "/assets/" + assetName
	log.Println("Serving asset: ", postAssetPath)
	c.File(postAssetPath)
	c.Status(200)
}

func Post(c *gin.Context) {
	postName := c.Param("post")
	post, err := getPost(postName)
	if err != nil {
		c.AbortWithError(500, err)
	}

	postComponent := templates.Unsafe(string(post.Content))
	templates.Post(post, postComponent).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func getPost(postName string) (*domain.Post, error) {
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

	post.Content = mdToHTML(content)
	return post, nil
}

func Index(c *gin.Context) {

	mds, err := filepath.Glob(pkg.POSTS_PATH + "**/*.md")
	if err != nil {
		c.AbortWithError(500, err)
	}

	metas, err := filepath.Glob(pkg.POSTS_PATH + "**/metadata.json")
	if err != nil {
		c.AbortWithError(500, err)
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
			c.AbortWithError(500, err)
		}

		post.Metadata = &domain.Metadata{}
		err = json.Unmarshal(metaBytes, post.Metadata)
		if err != nil {
			log.Println("Invalid metadata.json: ", meta, err)
		}
	}

	templates.Index(posts).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
