package main

import (
	"log"
	"os"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/service/minify"
)

func main() {
	if os.Getenv("MINIFY_ASSETS") != "1" {
		log.Println("[assets] MINIFY_ASSETS!=1; skipping")
		return
	}

	cfg := config.Config()

	if err := minify.NewWebMinifier(cfg).Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("[assets] done")
}
