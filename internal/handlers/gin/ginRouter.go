package ginhdl

import (
	"fmt"
	"log"
	"log/slog"
	"net/mail"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/contrib/gzip"
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
	router.Engine.Use(gzip.Gzip(gzip.DefaultCompression))
	router.registerRoutes()
	return router
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply the Cache-Control header to the static file
		c.Header("Cache-Control", "private, max-age=31536000")
		// Continue to the next middleware or handler
		c.Next()
	}
}

func (gr GinRouter) registerRoutes() {
	r := gr.Engine
	r.GET("/", gr.Index)
	r.NoRoute(gr.NoRoute)
	r.GET("/posts/:post", gr.Post)
	r.GET("/about", gr.About)
	r.POST("/search", gr.SearchPosts)
	r.GET("/subscribe", gr.Subscribe)
	r.POST("/subscribe", gr.SubscribeAdd)
	//Static files that should be served at root
	r.Use(staticCacheMiddleware())
	r.GET("/posts/:post/assets/:asset", gr.PostAssetAbs)
	r.GET("/posts/assets/:asset", gr.PostAsset)
	r.StaticFile("/output.css", fmt.Sprintf("%s/output.css", pkg.DIST_PATH))
	r.StaticFile("/site.webmanifest", fmt.Sprintf("%s/site.webmanifest", pkg.DIST_PATH))
	r.StaticFile("/favicon.ico", fmt.Sprintf("%s/favicon.ico", pkg.DIST_PATH))
	r.StaticFile("/favicon-32x32.png", fmt.Sprintf("%s/favicon-32x32.png", pkg.DIST_PATH))
	r.StaticFile("/favicon-16x16.png", fmt.Sprintf("%s/favicon-16x16.png", pkg.DIST_PATH))
	r.StaticFile("/apple-touch-icon.png", fmt.Sprintf("%s/apple-touch-icon.png", pkg.DIST_PATH))
	r.StaticFile("/android-chrome-512x512.png", fmt.Sprintf("%s/android-chrome-512x512.png", pkg.DIST_PATH))
	r.StaticFile("/android-chrome-192x192.png", fmt.Sprintf("%s/android-chrome-192x192.png", pkg.DIST_PATH))
	r.StaticFile("/robots.txt", fmt.Sprintf("%s/robots.txt", pkg.DIST_PATH))
	r.StaticFile("/about/about_profile.png", fmt.Sprintf("%s/about_profile.png", pkg.DIST_PATH))
}

func (gr GinRouter) About(c *gin.Context) {
	about := templates.About()
	frag := c.Request.URL.Query().Get("fragment") == "1"
	if !frag {
		about = templates.Base(state.BaseState{
			Title: "Guigoes - Guilherme de Castro",
			State: state.State{Language: getLanguage(c)},
			Body:  about,
		})
	}

	c.Header("HX-Replace-Url", "/about")
	c.Header("HX-Push-Url", "/about")
	c.Header("Content-Type", "text/html; charset=utf-8")
	about.Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func (gr GinRouter) SearchPosts(c *gin.Context) {
	search := c.Request.FormValue("search")
	if search == "" {

		ref := c.Request.Header.Get("Referer")
		if ref != "" {
			url, err := url.Parse(ref)
			if err != nil {
				slog.Error("Error parsing URL:", err)
				c.AbortWithError(500, err)
				return
			}

			post := strings.TrimPrefix(url.Path, "/posts/")
			if post != url.Path {
				gr.GetPostByName(post, true, c)
				return
			}
		}

		posts, err := gr.PostSrv.Posts(nil)
		if err != nil {
			slog.Error("Error retrieving posts:", err)
			c.AbortWithError(500, err)
			return
		}

		idxState := state.IndexState{
			State: state.State{Language: getLanguage(c)},
			Posts: posts,
		}

		templates.Index(idxState).Render(c.Request.Context(), c.Writer)
		c.Status(200)
		return
	}

	posts, err := gr.PostSrv.SearchPosts(search)
	if err != nil {
		slog.Error("Error searching posts:", err)
		c.AbortWithError(500, err)
		return
	}

	idxState := state.IndexState{
		State: state.State{Language: getLanguage(c)},
		Posts: posts,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Index(idxState).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func (gr GinRouter) Subscribe(c *gin.Context) {
	subComponent := templates.Subscribe("")
	modal := c.Query("modal")
	if modal == "1" {
		templates.Modal(subComponent).Render(c.Request.Context(), c.Writer)
		c.Status(200)
		return
	}

	subComponent.Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
func (gr GinRouter) SubscribeAdd(c *gin.Context) {
	email := c.Request.FormValue("email")
	addr, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("Bad email address %s\n", err)
		c.AbortWithError(400, err)
	}
	slog.Info("SUBSCRIBED:", email, addr)
	templates.SubscribeOk("").Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func (gr GinRouter) PostAsset(c *gin.Context) {
	ref := c.Request.Header.Get("Referer")
	if ref == "" {
		c.AbortWithStatus(404)
		return
	}

	url, err := url.Parse(ref)
	if err != nil {
		slog.Error("Error parsing URL:", err)
		c.AbortWithError(500, err)
		return
	}

	postName := filepath.Base(url.Path)
	assetName := c.Param("asset")

	assetPath, err := gr.PostSrv.GetPostAsset(postName, assetName)
	if err != nil {
		slog.Error("Error retrieving post asset:", err)
		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			c.AbortWithError(404, err)
			return
		}
		c.AbortWithError(500, err)
		return
	}

	slog.Debug("Serving asset: ", assetPath, "")
	c.File(assetPath)
	c.Status(200)
}

func (gr GinRouter) PostAssetAbs(c *gin.Context) {
	postName := c.Param("post")
	assetName := c.Param("asset")
	assetPath, err := gr.PostSrv.GetPostAsset(postName, assetName)
	if err != nil {
		slog.Error("Error retrieving post asset:", err)
		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			c.AbortWithError(404, err)
			return
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
	frag := c.Request.URL.Query().Get("fragment") == "1"
	gr.GetPostByName(postName, frag, c)
}

func (gr GinRouter) GetPostByName(postName string, frag bool, c *gin.Context) {
	post, err := gr.PostSrv.GetPost(postName)
	if err != nil {
		slog.Error("Error retrieving post:", err)

		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			gr.NoRoute(c)
			return
		}
		c.AbortWithError(500, err)
		return
	}

	c.Header("Last-Modified", post.Metadata.UpdatedAt.ToRfc7231String())
	postContent := templates.Unsafe(string(post.Content))
	postFragment := templates.Post(post, postContent)
	if frag {
		c.Header("HX-Replace-Url", post.Dir)
		c.Header("HX-Push-Url", post.Dir)
		postFragment.Render(c.Request.Context(), c.Writer)
		c.Status(200)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Base(state.BaseState{
		State: state.State{Language: getLanguage(c)},
		Title: post.Metadata.Title,
		Body:  postFragment,
		Post:  post,
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

	posts, err := gr.PostSrv.Posts(nil)
	if err != nil {
		slog.Error("Error retrieving posts:", err)
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

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Base(bs).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}

func (gr GinRouter) NoRoute(c *gin.Context) {

	idxState := state.IndexState{
		State: state.State{Language: getLanguage(c)},
		Posts: []*domain.Post{},
	}

	bs := state.BaseState{
		Title: "Guigoes - Home",
		Body:  templates.Index(idxState),
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Base(bs).Render(c.Request.Context(), c.Writer)
	c.Status(200)
}
