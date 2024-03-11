package main

import (
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
)

var postsService ports.PostService

func init() {
	pkg.LoadEnvFile()
	postsService = services.NewLocalPostService()
}

type FlatPost struct {
	Dir         string
	Name        string
	Title       string
	Author      string
	Description string
	Tags        []string
	PostName    string
	Content     string
	Type        string
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for portuguese text
	ptTextFieldMapping := bleve.NewTextFieldMapping()
	ptTextFieldMapping.Store = true
	ptTextFieldMapping.IncludeInAll = true
	ptTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Store = true
	keywordFieldMapping.IncludeInAll = true
	keywordFieldMapping.Analyzer = keyword.Name

	postMapping := bleve.NewDocumentMapping()

	// Metadata.Title
	postMapping.AddFieldMappingsAt("Dir", keywordFieldMapping)
	postMapping.AddFieldMappingsAt("Name", ptTextFieldMapping)
	postMapping.AddFieldMappingsAt("Title", ptTextFieldMapping)
	postMapping.AddFieldMappingsAt("Author", keywordFieldMapping)
	postMapping.AddFieldMappingsAt("Description", ptTextFieldMapping)
	postMapping.AddFieldMappingsAt("Tags", keywordFieldMapping)
	postMapping.AddFieldMappingsAt("PostName", ptTextFieldMapping)
	postMapping.AddFieldMappingsAt("Content", ptTextFieldMapping)
	postMapping.AddFieldMappingsAt("Type", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("post", postMapping)

	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func main() {
	// open the index
	blogIndex, err := bleve.Open("blog.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		blogIndex, err = bleve.New("blog.bleve", indexMapping)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	defer blogIndex.Close()

	posts, err := postsService.Posts(&ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1000)),
		Content: pkg.Ptr(ports.Markdown),
	})
	if err != nil {
		log.Fatalln(err)
	}

	batch := blogIndex.NewBatch()
	for _, post := range posts {
		batch.Index(post.Dir, FlatPost{
			Dir:         post.Dir,
			Name:        post.Name,
			Title:       post.Metadata.Title,
			Author:      post.Metadata.Author,
			Description: post.Metadata.Description,
			Tags:        post.Metadata.Tags,
			PostName:    post.Metadata.PostName,
			Content:     string(post.Content),
			Type:        "post",
		})
	}

	err = blogIndex.Batch(batch)
	if err != nil {
		log.Fatalln(err)
	}
	query := bleve.NewQueryStringQuery("HTMX")
	search := bleve.NewSearchRequest(query)
	searchResults, err := blogIndex.Search(search)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%v", searchResults)
}
