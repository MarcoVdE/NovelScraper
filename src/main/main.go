package main

import (
	"../crawler"
	"../rest_api"
	"image"
)

type NovelFetch struct {
	Title string
	URL   string
}

type NovelList struct {
	Title       string
	ChapterList []crawler.Chapter
	NovelImage  image.Image
}

func init() {
	var novels []NovelList
	//novels to get:
	var titles = []NovelFetch{
		{
			Title: "Sovereign Of The Karmic System",
			URL:   "sovereign-of-the-karmic-system",
		},
	}

	for _, each := range titles {
		chapter, novelImage := crawler.WuxiaWorldCrawler(each.URL)
		novels = append(novels, NovelList{each.Title, chapter, novelImage})
	}

	//TODO: Create API for fetching the books. Use JSON Encode.
	//TODO: Create integration with: https://github.com/languagetool-org/languagetool
}

func main() {
	rest_api.RestAPIEntry()

}
