package main

import (
	"flag"
	"net/http"
	"path/filepath"

	stdhdl "github.com/guilycst/guigoes/internal/handlers/std"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
	"github.com/otiai10/copy"
)

var postsService ports.PostService

func init() {
	envfile := flag.String("envfile", ".env", "path to env file")
	flag.Parse()
	pkg.LoadEnvFile(*envfile)
	postsService = services.NewLocalPostService()
}

func main() {
	idxTmp := "/tmp/" + filepath.Base(pkg.BLEVE_IDX_PATH) + "/"
	copy.Copy(pkg.BLEVE_IDX_PATH, idxTmp)
	pkg.BLEVE_IDX_PATH = idxTmp

	// mux := http.NewServeMux()

	// // Register the routes and handlers
	// mux.Handle("/", &homeHandler{})

	// Run the server
	http.ListenAndServe(":8080", stdhdl.NewStandardRouter(postsService))
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
