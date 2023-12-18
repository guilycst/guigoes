package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/pkg"
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
	//router.Engine.Use(gzip.Gzip(gzip.DefaultCompression))
	router.registerRoutes()
	return router
}

func (gr GinRouter) registerRoutes() {
	r := gr.Engine
	r.GET("/", gr.Index)
	r.GET("/posts/:post", gr.Post)
	r.GET("/posts/assets/:asset", gr.PostAsset)
	//Static files that should be served at root
	r.StaticFile("/output.css", fmt.Sprintf("%s/output.css", pkg.DIST_PATH))
	r.StaticFile("/site.webmanifest", fmt.Sprintf("%s/site.webmanifest", pkg.DIST_PATH))
	r.StaticFile("/favicon.ico", fmt.Sprintf("%s/favicon.ico", pkg.DIST_PATH))
	r.StaticFile("/favicon-32x32.png", fmt.Sprintf("%s/favicon-32x32.png", pkg.DIST_PATH))
	r.StaticFile("/favicon-16x16.png", fmt.Sprintf("%s/favicon-16x16.png", pkg.DIST_PATH))
	r.StaticFile("/apple-touch-icon.png", fmt.Sprintf("%s/apple-touch-icon.png", pkg.DIST_PATH))
	r.StaticFile("/android-chrome-512x512.png", fmt.Sprintf("%s/android-chrome-512x512.png", pkg.DIST_PATH))
	r.StaticFile("/android-chrome-192x192.png", fmt.Sprintf("%s/android-chrome-192x192.png", pkg.DIST_PATH))
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

	c.Header("Last-Modified", post.UpdatedAt.ToRfc7231String())
	frag := c.Request.URL.Query().Get("fragment") == "1"
	postContent := templates.Unsafe(string(post.Content))
	postFragment := templates.Post(post, postContent)
	if frag {
		c.Header("HX-Replace-Url", post.Dir)
		postFragment.Render(c.Request.Context(), c.Writer)
		c.Status(200)
		return
	}

	templates.Base(state.BaseState{
		State: state.State{Language: getLanguage(c)},
		Title: post.Metadata.Title,
		Body:  postFragment,
	}).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func getLanguage(c *gin.Context) string {
	header := c.Request.Header.Get("Accept-Language")
	if header == "" {
		return "en"
	}

	langs := strings.Split(header, ";")
	if len(langs) == 0 {
		return "en"
	}

	for _, lang := range langs {
		return strings.Split(lang, ",")[1]
	}

	return "en"
}

func (gr GinRouter) Index(c *gin.Context) {

	posts, err := gr.PostSrv.Index()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	idxState := state.IndexState{
		State: state.State{Language: getLanguage(c)},
		Posts: posts,
	}

	bs := state.BaseState{
		Title: "Guigoes - Home",
		Body:  templates.Index(idxState),
	}
	templates.Base(bs).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
