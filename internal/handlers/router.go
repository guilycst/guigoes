package handlers

import (
	"errors"
	"log"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/web/templates"
	"github.com/guilycst/guigoes/web/templates/state"
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
	router.Engine.Use(gzip.Gzip(gzip.DefaultCompression))
	router.registerRoutes()
	return router
}

func (gr GinRouter) registerRoutes() {
	r := gr.Engine
	r.GET("/", gr.Index)
	r.GET("/posts/:post", gr.Post)
	r.GET("/posts/assets/:asset", gr.PostAsset)
	//Static files that should be served at root
	r.StaticFile("/output.css", "./web/dist/output.css")
	r.StaticFile("/site.webmanifest", "./web/dist/site.webmanifest")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	r.StaticFile("/favicon-32x32.png", "./web/dist/favicon-32x32.png")
	r.StaticFile("/favicon-16x16.png", "./web/dist/favicon-16x16.png")
	r.StaticFile("/apple-touch-icon.png", "./web/dist/apple-touch-icon.png")
	r.StaticFile("/android-chrome-512x512.png", "./web/dist/android-chrome-512x512.png")
	r.StaticFile("/android-chrome-192x192.png", "./web/dist/android-chrome-192x192.png")
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

	bs := state.BaseState{
		Title: "Guigoes - Home",
		Body:  templates.Index(idx),
	}
	templates.Base(bs).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
