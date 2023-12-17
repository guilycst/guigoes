package handlers

import (
	"errors"
	"log"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/web/templates"
)

type GinRouter struct {
	Engine  *gin.Engine
	PostSrv ports.PostService
}

func NewGinRouter(ps ports.PostService) *GinRouter {
	router := &GinRouter{
		Engine:  gin.Default(),
		PostSrv: ps,
	}
	router.registerRoutes()
	return router
}

func (gr GinRouter) registerRoutes() {
	r := gr.Engine
	r.GET("/", gr.Index)
	r.GET("/posts/:post", gr.Post)
	r.GET("/posts/assets/:asset", gr.PostAsset)
}

func (gr GinRouter) PostAsset(c *gin.Context) {
	ref := c.Request.Header.Get("Referer")
	if ref == "" {
		c.AbortWithStatus(404)
		return
	}

	url, err := url.Parse(ref)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	postName := filepath.Base(url.Path)
	assetName := c.Param("asset")

	assetPath, err := gr.PostSrv.GetPostAsset(postName, assetName)
	if err != nil {
		if errors.Is(err, &domain.AssetNotFoundError{}) {
			c.AbortWithError(404, err)
		}
		c.AbortWithError(500, err)
		return
	}

	log.Println("Serving asset: ", assetPath)
	c.File(assetPath)
	c.Status(200)
}

func (gr GinRouter) Post(c *gin.Context) {
	postName := c.Param("post")
	post, err := gr.PostSrv.GetPost(postName)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	postComponent := templates.Unsafe(string(post.Content))
	templates.Post(post, postComponent).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func (gr GinRouter) Index(c *gin.Context) {
	idx, err := gr.PostSrv.Index()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	templates.Index(idx).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
