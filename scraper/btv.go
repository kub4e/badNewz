package main

import (
	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

func scrapteBTV(depth int, articleList *[]article) {
	url := "https://btvnovinite.bg/bulgaria/?page="

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"),
	)

	c.OnHTML(".news-article-inline", func(e *colly.HTMLElement) {

		a := article{
			Title:       e.ChildText(".title"),
			Description: e.ChildText(".text"),
			Url:         "https://btvnovinite.bg" + e.ChildAttr("a", "href"),
			Date:        e.ChildText(".date"),
			Img:         e.ChildAttr("img", "src"),
			Source:      "BTV",
		}
		a.Content = getContent(a.Url, ".article-body")
		*articleList = append(*articleList, a)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i <= depth; i++ {
		c.Visit(url + strconv.Itoa(i))
	}
}
