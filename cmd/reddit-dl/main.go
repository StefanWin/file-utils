package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/StefanWin/file-util/v2/pkg"
)

func main() {

	var subreddit string
	flag.StringVar(&subreddit, "subreddit", "", "")
	var flair string
	flag.StringVar(&flair, "flair", "", "")
	flag.Parse()

	if subreddit == "" {
		log.Fatal(fmt.Errorf("no subreddit"))
	}

	var images []string
	filename := ""
	if flair == "" {
		urls, err := pkg.GetSubRedditImageURLS(subreddit)
		if err != nil {
			log.Fatal(err)
		}
		images = urls
		filename = fmt.Sprintf("%s.txt", subreddit)
	} else {
		urls, err := pkg.GetSubRedditImageURLSFlair(subreddit, flair)
		if err != nil {
			log.Fatal(err)
		}
		images = urls
		filename = fmt.Sprintf("%s_%s.txt", subreddit, flair)
	}
	log.Printf("found %d images\n", len(images))
	txt := strings.Join(images, "\n")
	if err := os.WriteFile(filename, []byte(txt), 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote urls to %s\n", filename)
}
