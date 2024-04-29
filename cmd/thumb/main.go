package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"path"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

var postsService ports.PostService

func init() {
	envfile := flag.String("envfile", ".env", "path to env file")
	flag.Parse()
	pkg.LoadEnvFile(*envfile)
	postsService = services.NewLocalPostService()
}

func main() {
	posts, err := postsService.Posts(nil)
	if err != nil {
		panic(err)
	}

	for _, v := range posts {
		coverImgPath := path.Join(".", v.Dir, v.Metadata.Cover)
		f, err := os.Open(coverImgPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		src, _, err := image.Decode(f)
		if err != nil {
			panic(err)
		}

		dst := image.NewRGBA(image.Rect(0, 0, 300, 198))
		draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

		thumbImgPath := path.Base(coverImgPath)
		thumbImgPath = path.Join(".", v.Dir, "assets", "thumb_"+thumbImgPath)
		tf, err := os.Create(thumbImgPath)
		if err != nil {
			panic(err)
		}
		defer tf.Close()

		png.Encode(tf, dst)
	}
}
