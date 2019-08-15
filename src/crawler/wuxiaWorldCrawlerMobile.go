package crawler

import (
	"../models"
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

func cleanContentByStringSlice(content string, cleanup []string) string {
	for _, each := range cleanup {
		content = strings.ReplaceAll(content, each, "\n\n")
	}
	return content
}

func cleanContent(content string) string {
	//remove first div with addsense
	//removing the Google Add at the beginning.
	parts := strings.Split(content, `</div>`)
	fmt.Printf("%q\n", parts)

	if len(parts) > 1 {
		content = parts[1]
	} else {
		content = parts[0]
	}
	//remove amp auto adds
	strings.ReplaceAll(content, `<amp-auto-ads type="adsense"
	data-ad-client="ca-pub-2853920792116568">
	</amp-auto-ads>
	<script>app2()</script>`, "<br>")

	//remove chapter mid
	strings.ReplaceAll(content, `<script>ChapterMid();</script>`, "<br>")

	//remove end div.
	parts = strings.Split(content, "<div>")
	fmt.Printf("%q\n", parts)

	if len(parts) > 1 {
		content = parts[0]
	} else {
		content = parts[1]
	}

	return content
}

//this is based on the CourseRA crawler.
func WuxiaWorldMobileCrawler(novelTitle string) ([]models.Chapter, image.Image) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains("wuxiaworld.co", "www.wuxiaworld.co", "m.wuxiaworld.co"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./wuxiaworld_mobile_cache"),
	)
	// Create another collector to scrape course details
	detailCollector := c.Clone()
	chapters := make([]models.Chapter, 0, 200)

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`div[id='chapterlist']>p>a[href$='.html']`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains chapter
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visiting link: " + courseURL)
		if strings.Index(courseURL, ".html") != -1 {
			fmt.Println("Visiting link: " + courseURL)
			//TODO: Maybe pass title here? It's in the content of the a tag, or add url to chapter and then associate two slices.
			_ = detailCollector.Visit(courseURL)
		}
	})

	// Extract details of the course
	detailCollector.OnHTML(`body[id='read']`, func(e *colly.HTMLElement) {
		fmt.Printf("Title found at URL: %s", e.Request.URL)
		title := e.ChildText("header[id='top']>span.title")
		if title == "" {
			fmt.Println("No title found", e.Request.URL)
		}

		chapter := models.Chapter{
			Title:   title,
			URL:     e.Request.URL.String(),
			Chapter: getMobileChapterFromTitle(title),
			Content: cleanContentByStringSlice(
				e.ChildText("div[id=chaptercontent]"),
				[]string{
					//`<amp-auto-ads type="adsense"
					//data-ad-client="ca-pub-2853920792116568">
					//</amp-auto-ads>`,
					//"<script>app2()</script>",
					"app2()",
					//"<script>ChapterMid();</script>",
					"ChapterMid();",
					e.ChildText("div[id=chaptercontent] div"),
				},
			),
			//Content: cleanContent(e.ChildText("div[id=chaptercontent]")),
		}
		//filter out empty content like preview ones.
		if chapter.Content != "" {
			chapters = append(chapters, chapter)
		}
	})

	_ = c.Visit("https://m.wuxiaworld.co/" + novelTitle + "/all.html")

	//TODO: Scrape base page, can get chapter information from there, and check whether there is a new chapter since update time, reducing load on site.

	//All images follow rule: https://www.wuxiaworld.co/BookFiles/BookImages/sovereign-of-the-karmic-system.jpg Note capitalization doesn't matter.
	return chapters, getMobileImage("https://www.wuxiaworld.co/BookFiles/BookImages/" + novelTitle + ".jpg")
}
