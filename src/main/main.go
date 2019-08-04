package main

type NovelList struct {
	Title string
	ChapterList []Chapter
}

func main() {
	//novels to get:
	var novels []NovelList
	novels = append(novels,
		NovelList{Title: "Sovereign Of The Karmic System", ChapterList: wuxiaWorldCrawler("sovereign-of-the-karmic-system")})

	//TODO: Create API for fetching the books. Use JSON Encode.
	//TODO: Create
}

