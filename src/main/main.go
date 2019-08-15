package main

import (
	"../crawler"
	model "../models"
	"../rest_api"
)

//var novels []model.NovelList
var mobileNovels []model.NovelList

func init() {
	//novels to get:
	var titles = []model.NovelFetch{
		{
			Title: "Sovereign Of The Karmic System",
			URL:   "sovereign-of-the-karmic-system",
		},
	}
	//
	//for _, each := range titles {
	//	chapter, novelImage := crawler.WuxiaWorldCrawler(each.URL)
	//	novels = append(novels, model.NovelList{
	//		Title: each.Title,
	//		ChapterList: chapter,
	//		NovelImage: novelImage,
	//	})
	//}
	//
	//titles = []model.NovelFetch {
	//	{
	//		Title: "Sovereign Of The Karmic System",
	//		URL: "sovereign-of-the-karmic-system",
	//	},
	//}

	for _, each := range titles {
		chapter, novelImage := crawler.WuxiaWorldMobileCrawler(each.URL)
		mobileNovels = append(mobileNovels, model.NovelList{
			Title:       each.Title,
			ChapterList: chapter,
			NovelImage:  novelImage,
		})
	}

	//TODO: Create integration with: https://github.com/languagetool-org/languagetool
	//maybe do this client side to get over 100 request a day limit? User can then send it back, if x amount of user send it's good.
	//or allow instantly if user certain privileges (e.g. me) after authentication implemented.
}

func main() {
	//rest_api.RestAPIEntry(novels)
	rest_api.RestAPIEntry(mobileNovels)
}
