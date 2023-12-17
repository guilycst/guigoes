package main

import (
	"github.com/guilycst/guigoes/internal/handlers"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
)

var postsService ports.PostService

func init() {
	pkg.LoadEnvFile()
	postsService = services.NewLocalPostService()
}

func main() {
	r := handlers.NewGinRouter(postsService)
	r.Engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
