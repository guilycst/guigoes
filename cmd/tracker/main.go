package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/guilycst/guigoes/internal/ports"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
)

const postCommitsDatetimeHashCmd = `log --format=%%aI_%%H -- %s`

var postsService ports.PostService

func init() {
	pkg.LoadEnvFile()
	postsService = services.NewLocalPostService()
}

func main() {
	idx, err := postsService.Posts(&ports.PostsOptions{
		FetchN:  pkg.Ptr(uint64(1000)),
		Content: pkg.Ptr(ports.None),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range idx {
		bodyDir := "." + p.Dir
		bodyPath := bodyDir + "/body.md"
		tracker := bodyDir + "/git-log.track"
		cmdParams := fmt.Sprintf(postCommitsDatetimeHashCmd, bodyPath)
		cmd := exec.Command("git", strings.Split(cmdParams, " ")...)
		var out strings.Builder
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(tracker)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write([]byte(out.String()))

		if err != nil {
			log.Fatal(err)
		}
	}
}
