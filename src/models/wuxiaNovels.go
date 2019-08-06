package models

import (
	"image"
)

type NovelFetch struct {
	Title string
	URL   string
}

type NovelList struct {
	Title       string      `json:"Title"`
	ChapterList []Chapter   `json:"ChapterList"`
	NovelImage  image.Image `json:"NovelImage"`
}

type Chapter struct {
	Title   string `json:"Title"` //section#content>div.my_container>div.content>div.content_left>div.manga_view_name>h1.
	URL     string `json:"URL"`
	Chapter int    `json:"Chapter"` //extract from title as always Chapter space number.
	Content string `json:"Content"` //div#content //Note <br><br> is used for next line. Extracted converts to \n
}
