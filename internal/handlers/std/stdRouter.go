package std

import (
	"log"
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/guilycst/guigoes/internal/core/domain"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/pkg"
	"github.com/guilycst/guigoes/pkg/middleware"
	"github.com/guilycst/guigoes/web/templates"
	"github.com/guilycst/guigoes/web/templates/state"
)

type StandardRouter struct {
	PostSrv ports.PostService
	handler *http.ServeMux
}

func NewStandardRouter(ps ports.PostService) *StandardRouter {
	mux := http.NewServeMux()
	r := &StandardRouter{
		PostSrv: ps,
		handler: mux,
	}

	r.registerRoutes()

	return r
}

type MiddlewareMux struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

func newMiddlewareMux(mux *http.ServeMux) *MiddlewareMux {
	mm := MiddlewareMux{
		mux:        mux,
		middleware: []func(http.Handler) http.Handler{},
	}
	return &mm
}

func (mm *MiddlewareMux) Use(middleware ...func(http.Handler) http.Handler) {
	mm.middleware = append(mm.middleware, middleware...)
}

func (mm *MiddlewareMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var h http.Handler = http.HandlerFunc(handler)
	for _, m := range mm.middleware {
		h = m(h)
	}

	mm.mux.Handle(pattern, h)
}

func (mm *MiddlewareMux) Clone() *MiddlewareMux {
	nm := newMiddlewareMux(mm.mux)
	nm.middleware = append(nm.middleware, mm.middleware...)
	return nm
}

func (sr *StandardRouter) registerRoutes() {
	mux := newMiddlewareMux(sr.handler)
	mux.Use(middleware.PanicRecover, middleware.Gzip)

	mux.HandleFunc("GET /", sr.Index)
	mux.HandleFunc("GET /posts/{post}", sr.Post)
	mux.HandleFunc("GET /about", sr.About)
	mux.HandleFunc("POST /search", sr.SearchPosts)
	mux.HandleFunc("GET /subscribe", sr.Subscribe)
	mux.HandleFunc("POST /subscribe", sr.SubscribeAdd)

	//those should have a cache control header
	mux = mux.Clone()
	mux.Use(middleware.CacheControl)

	mux.HandleFunc("GET /posts/{post}/assets/{asset}", sr.PostAssetAbs)
	mux.HandleFunc("GET /posts/assets/{asset}", sr.PostAsset)
	mux.HandleFunc("GET /output.css", sr.StaticFileAtomic("text/css; charset=utf-8"))
	mux.HandleFunc("GET /site.webmanifest", sr.StaticFileAtomic("text/plain; charset=utf-8"))
	mux.HandleFunc("GET /favicon.ico", sr.StaticFile)
	mux.HandleFunc("GET /favicon-32x32.png", sr.StaticFile)
	mux.HandleFunc("GET /favicon-16x16.png", sr.StaticFile)
	mux.HandleFunc("GET /apple-touch-icon.png", sr.StaticFile)
	mux.HandleFunc("GET /android-chrome-512x512.png", sr.StaticFile)
	mux.HandleFunc("GET /android-chrome-192x192.png", sr.StaticFile)
	mux.HandleFunc("GET /robots.txt", sr.StaticFileAtomic("text/plain; charset=utf-8"))
	mux.HandleFunc("GET /about/about_profile.png", sr.StaticFile)

}

func (sr *StandardRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sr.handler.ServeHTTP(w, r)
}

func (sr *StandardRouter) StaticFile(w http.ResponseWriter, r *http.Request) {
	assetPath := r.URL.Path
	// Make sure filePath is clean to avoid path traversal attacks
	assetPath = filepath.Clean(assetPath)
	// Serve the file from the specified directory
	http.ServeFile(w, r, filepath.Join(pkg.DIST_PATH, filepath.Base(assetPath)))
}

func (sr *StandardRouter) StaticFileAtomic(contentType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		assetPath := r.URL.Path
		// Make sure filePath is clean to avoid path traversal attacks
		assetPath = filepath.Clean(assetPath)
		assetPath = filepath.Join(pkg.DIST_PATH, filepath.Base(assetPath))
		f, err := os.ReadFile(assetPath)
		if err != nil {
			slog.Error("Error reading file from disk", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		_, err = w.Write(f)
		if err != nil {
			slog.Error("Error writing file to response", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (sr *StandardRouter) About(w http.ResponseWriter, r *http.Request) {
	about := templates.About()
	frag := r.URL.Query().Get("fragment") == "1"
	if !frag {
		about = templates.Base(state.BaseState{
			Title: "Guigoes - Guilherme de Castro",
			State: state.State{Language: getLanguage(r)},
			Body:  about,
		})
	}

	w.Header().Set("HX-Replace-Url", "/about")
	w.Header().Set("HX-Push-Url", "/about")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	about.Render(r.Context(), w)
}

func (sr *StandardRouter) SearchPosts(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	if search == "" {
		ref := r.Header.Get("Referer")
		if ref != "" {
			url, err := url.Parse(ref)
			if err != nil {
				slog.Error("Error parsing URL:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := strings.TrimPrefix(url.Path, "/posts/")
			if post != url.Path {
				sr.GetPostByName(post, true, w, r)
				return
			}
		}

		posts, err := sr.PostSrv.Posts(nil)
		if err != nil {
			slog.Error("Error getting posts:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idxState := state.IndexState{
			State: state.State{Language: getLanguage(r)},
			Posts: posts,
		}

		templates.Index(idxState).Render(r.Context(), w)
		return
	}

	posts, err := sr.PostSrv.SearchPosts(search)
	if err != nil {
		slog.Error("Error searching posts:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idxState := state.IndexState{
		State: state.State{Language: getLanguage(r)},
		Posts: posts,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Index(idxState).Render(r.Context(), w)
}

func (sr *StandardRouter) Subscribe(w http.ResponseWriter, r *http.Request) {
	subComponent := templates.Subscribe("")
	modal := r.URL.Query().Get("modal")
	if modal == "1" {
		templates.Modal(subComponent).Render(r.Context(), w)
		return
	}

	subComponent.Render(r.Context(), w)
}

func (sr *StandardRouter) SubscribeAdd(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	addr, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("Bad email address %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("SUBSCRIBED:", email, addr)
	templates.SubscribeOk("").Render(r.Context(), w)
}

func (sr StandardRouter) PostAsset(w http.ResponseWriter, r *http.Request) {
	ref := r.Header.Get("Referer")
	if ref == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := url.Parse(ref)
	if err != nil {
		slog.Error("Error parsing URL:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postName := filepath.Base(url.Path)
	assetName := r.PathValue("asset")

	assetPath, err := sr.PostSrv.GetPostAsset(postName, assetName)
	if err != nil {
		slog.Error("Error getting post asset:", err)
		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.ReadFile(assetPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(f)
}

func (sr *StandardRouter) PostAssetAbs(w http.ResponseWriter, r *http.Request) {
	postName := r.PathValue("post")
	assetName := r.PathValue("asset")
	assetPath, err := sr.PostSrv.GetPostAsset(postName, assetName)
	if err != nil {
		slog.Error("Error getting post asset at path", assetPath, err)
		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("Serving asset", assetPath, "")
	http.ServeFile(w, r, assetPath)
}

func (sr *StandardRouter) Post(w http.ResponseWriter, r *http.Request) {
	postName := r.PathValue("post")
	frag := r.URL.Query().Get("fragment") == "1"
	sr.GetPostByName(postName, frag, w, r)
}

func (sr *StandardRouter) GetPostByName(postName string, frag bool, w http.ResponseWriter, r *http.Request) {
	post, err := sr.PostSrv.GetPost(postName)
	if err != nil {
		slog.Error("Error getting post:", err)

		if _, ok := err.(*domain.FSResourceNotFoundError); ok {
			sr.NoRoute(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Last-Modified", post.Metadata.UpdatedAt.ToRfc7231String())
	postContent := templates.Unsafe(string(post.Content))
	postFragment := templates.Post(post, postContent)
	if frag {
		w.Header().Set("HX-Replace-Url", post.Dir)
		w.Header().Set("HX-Push-Url", post.Dir)
		postFragment.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Base(state.BaseState{
		State: state.State{Language: getLanguage(r)},
		Title: post.Metadata.Title,
		Body:  postFragment,
		Post:  post,
	}).Render(r.Context(), w)
}

func (sr *StandardRouter) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := sr.PostSrv.Posts(nil)
	if err != nil {
		slog.Error("Error getting posts:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idxState := state.IndexState{
		State: state.State{Language: getLanguage(r)},
		Posts: posts,
	}

	bs := state.BaseState{
		Title: "Guigoes - Home",
		Body:  templates.Index(idxState),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Base(bs).Render(r.Context(), w)
}

func (sr *StandardRouter) NoRoute(w http.ResponseWriter, r *http.Request) {
	idxState := state.IndexState{
		State: state.State{Language: getLanguage(r)},
		Posts: []*domain.Post{},
	}

	bs := state.BaseState{
		Title: "Guigoes - Home",
		Body:  templates.Index(idxState),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Base(bs).Render(r.Context(), w)
}

func getLanguage(r *http.Request) string {
	header := r.Header.Get("Accept-Language")
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
