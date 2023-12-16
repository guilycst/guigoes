package handlers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/web/templates"
)

func NewGin() *gin.Engine {
	r := gin.Default()
	r.GET("/", Index)
	return r
}

func Index(c *gin.Context) {
	var root = "./posts/"

	mds, err := filepath.Glob(root + "**/*.md")
	if err != nil {
		c.AbortWithError(500, err)
	}

	metas, err := filepath.Glob(root + "**/metadata.json")
	if err != nil {
		c.AbortWithError(500, err)
	}

	posts := make(map[string]*domain.Post)
	for _, md := range mds {
		dir := filepath.Dir(md)
		posts[dir] = &domain.Post{Dir: filepath.Base(md)}
	}

	for _, meta := range metas {
		dir := filepath.Dir(meta)
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
