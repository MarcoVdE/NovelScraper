package crawler

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"../models"
	"github.com/gocolly/colly"
)

//Get image off of: https://m.wuxiaworld.co/Sovereign-of-the-Karmic-System/
//Get content off of: https://m.wuxiaworld.co/Sovereign-of-the-Karmic-System/all.html //note the rule of all.html to content
func getMobileChapterFromTitle(title string) int {
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

func getMobileImage(url string) image.Image {
	//get cover image
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	img, _, err := image.Decode(response.Body)
	_ = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return img
}

//this is based on the CourseRA crawler.
func WuxiaWorldMobileCrawler(novelTitle string) ([]models.Chapter, image.Image) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains("wuxiaworld.co", "wuxiaworld.co"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./wuxiaworld_mobile_cache"),
	)

	//Novel cover image.
	var novelImage image.Image

	// Create another collector to scrape course details
	detailCollector := c.Clone()
	chapters := make([]models.Chapter, 0, 200)

	//On every a element which has href attribute call callback
	// c.OnHTML("a[href~=.html']", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// start scraping the page under the link found
	// 	_ = e.Request.Visit(link)
	// })

	//TODO: Make sure this goes to /all.html
	// On every a HTML element which has name attribute call callback
	c.OnHTML(`div[id=chapterlist] > p > a[href]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains chapter
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visiting link: " + courseURL)
		if strings.Index(courseURL, ".html") != -1 {
			fmt.Println("Visiting link: " + courseURL)
			//TODO: Maybe pass title here? It's in the content of the a tag, or add url to chapter and then associate two slices.
			_ = detailCollector.Visit(courseURL)
		}
	})

	// TODO: Add image support. this file:122
	// c.OnHTML(`div.manga_info_img`, func(e *colly.HTMLElement) {
	// 	novelImage = getMobileImage(e.ChildAttr("img.img-responsive", "src"))
	// })

	// Extract details of the course
	detailCollector.OnHTML(`div[id=chaptercontent]`, func(e *colly.HTMLElement) {
		log.Println("Chapter found", e.Request.URL)
		title := e.ChildText("div.content h1")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		chapter := models.Chapter{
			Title:   title, //TODO: Title comes from elsewhere, this function needs to be split into two.
			URL:     e.Request.URL.String(),
			Chapter: getMobileChapterFromTitle(title),
			Content: strings.TrimSuffix(e.ChildText("div#content"), "chaptererror();"), //TODO: Need to remove the beginning div as it's a Google Add.
		}

		//filter out empty content like preview ones.
		if chapter.Content != "" {
			chapters = append(chapters, chapter)
		}
	})

	_ = c.Visit("https://m.wuxiaworld.co/" + novelTitle + "/all.html")

	//only get the image form here? TODO: Make own function
	_ = c.Visit("https://m.wuxiaworld.co/" + novelTitle + "/")

	return chapters, novelImage
}
