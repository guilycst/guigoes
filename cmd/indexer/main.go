package main

import (
	"log"

	"github.com/blevesearch/bleve/v2"
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
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("posts.bleve", mapping)
	if err != nil {
		log.Fatalln(err)
	}

	post, err := postsService.GetPost("home")
	if err != nil {
		log.Fatalln(err)
	}

	err = index.Index(post.Dir, "There are six levels of headings. They correspond with the six levels of HTML headings. You've probably noticed them already in the page. Each level down uses one more hash character.")
	if err != nil {
		log.Fatalln(err)
	}

	query := bleve.NewMatchQuery("levels")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%v", searchResults)
}
