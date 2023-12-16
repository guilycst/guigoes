package main

import (
	"github.com/guilycst/guigoes/internal/handlers"
)

func main() {
	r := handlers.NewGin()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
