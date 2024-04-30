package main

import (
	"flag"
	"path/filepath"

	"github.com/guilycst/guigoes/internal/handlers"
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
	r := handlers.NewGinRouter(postsService)
	r.Engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
