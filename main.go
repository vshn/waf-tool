package main

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	fmt.Println("Ciao Belli")
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	fmt.Println("clone")
	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/src-d/go-git",
		Progress: os.Stdout,
	})
	log.Println(err)
}
