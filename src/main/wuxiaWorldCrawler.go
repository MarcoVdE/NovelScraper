package main

import (
	//"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	//"os"
	"regexp"
	"strconv"
	"strings"
)

type Chapter struct {
	Title string //section#content > div.my_container > div.content > div.content_left > div.manga_view_name > h1. Extract chapter
	URL string
	Chapter int //extract from title as always Chapter space number.
	Content string //div#content //Note <br><br> is used for next line.
}

func getChapterFromTitle(title string) int {
	//title := "Sovereign of the Karmic System Chapter 187: The Hellblazer Company" //note space between results,
	re := regexp.MustCompile("Chapter ?[0-9]+") //note ? means optional space before it.
	reNum := regexp.MustCompile("[0-9]+")
	//Then get number out of it only.
	//Note [0] is because regexp returns array of results, we only need first.

	if !(len(re.FindAllString(title, -1)) > 0) {
		return -1 //this rule happens because the site does not use 0 as first intro chapter but -1.
	}

	chapter, err := strconv.Atoi(reNum.FindAllString(re.FindAllString(title, -1)[0], -1)[0])
	if err != nil {
		fmt.Print(err)
	}
	return chapter
}

//this is based on the CourseRA crawler.
func wuxiaWorldCrawler(novelTitle string) []Chapter {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains("wuxiaworld.world", "www.wuxiaworld.world"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./wuxiaworld_cache"),
	)

	// Create another collector to scrape course details
	detailCollector := c.Clone()
	chapters := make([]Chapter, 0, 200)

	//On every a element which has href attribute call callback
	c.OnHTML("a[href~='" + novelTitle + "/chapter-']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// start scraping the page under the link found
		_ = e.Request.Visit(link)
	})

	// Before making a request print "Visiting ..."
	//c.OnRequest(func(r *colly.Request) {
	//	log.Println("visiting", r.URL.String())
	//})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains chapter
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visiting link: " + courseURL)
		if strings.Index(courseURL, novelTitle + "/chapter-") != -1 {
			fmt.Println("Visiting link: " + courseURL)
			_ = detailCollector.Visit(courseURL)
		}
	})

	// Extract details of the course
	detailCollector.OnHTML(`section[id=content]`, func(e *colly.HTMLElement) {
		log.Println("Chapter found", e.Request.URL)
		title := e.ChildText("div.content h1")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		chapter := Chapter {
			Title:       	title,
			URL:         	e.Request.URL.String(),
			Chapter: 		getChapterFromTitle(title),
			Content: 		strings.TrimSuffix(e.ChildText("div#content"), "chaptererror();"),
		}

		//filter out empty content like preview ones.
		if chapter.Content != "" {
			chapters = append(chapters, chapter)
		}
	})

	_ = c.Visit("https://wuxiaworld.world/" + novelTitle + "/")
	//
	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", "  ")
	//
	//// Dump json to the standard output
	//_ = enc.Encode(chapters)
	return chapters
}
