package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Chapter struct {
	Title   string `json:"Title"` //section#content>div.my_container>div.content>div.content_left>div.manga_view_name>h1.
	URL     string `json:"URL"`
	Chapter int    `json:"Chapter"` //extract from title as always Chapter space number.
	Content string `json:"Content"` //div#content //Note <br><br> is used for next line. Extracted converts to \n
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

func getImage(url string) image.Image {
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
func WuxiaWorldCrawler(novelTitle string) ([]Chapter, image.Image) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains("wuxiaworld.world", "www.wuxiaworld.world"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./wuxiaworld_cache"),
	)

	//Novel cover image.
	var novelImage image.Image

	// Create another collector to scrape course details
	detailCollector := c.Clone()
	chapters := make([]Chapter, 0, 200)

	//On every a element which has href attribute call callback
	c.OnHTML("a[href~='"+novelTitle+"/chapter-']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// start scraping the page under the link found
		_ = e.Request.Visit(link)
	})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains chapter
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visiting link: " + courseURL)
		if strings.Index(courseURL, novelTitle+"/chapter-") != -1 {
			fmt.Println("Visiting link: " + courseURL)
			_ = detailCollector.Visit(courseURL)
		}
	})

	c.OnHTML(`div.manga_info_img`, func(e *colly.HTMLElement) {
		novelImage = getImage(e.ChildAttr("img.img-responsive", "src"))
	})

	// Extract details of the course
	detailCollector.OnHTML(`section[id=content]`, func(e *colly.HTMLElement) {
		log.Println("Chapter found", e.Request.URL)
		title := e.ChildText("div.content h1")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		chapter := Chapter{
			Title:   title,
			URL:     e.Request.URL.String(),
			Chapter: getChapterFromTitle(title),
			Content: strings.TrimSuffix(e.ChildText("div#content"), "chaptererror();"),
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
	return chapters, novelImage
}
