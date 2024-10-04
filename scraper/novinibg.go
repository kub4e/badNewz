package main

import (
	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

func scrapeNoviniBG(depth int, articleList *[]article, topic string) {
	urls := make(map[string]string)
	urls["politics"] = "https://novini.bg/bylgariya/politika?page="
	urls["economy"] = "https://novini.bg/biznes?page="

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"),
	)

	c.OnHTML("article.js-content", func(e *colly.HTMLElement) {

		a := article{
			Title:  e.ChildAttr("img", "alt"),
			Url:    e.ChildAttr("a", "href"),
			Date:   e.ChildText("span"),
			Img:    e.ChildAttr("img", "src"),
			Topic:  topic,
			Source: "NoviniBG",
		}
		a.Content = getContent(a.Url, ".openArticle__content")
		*articleList = append(*articleList, a)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	for i := 2; i <= depth; i++ {
		c.Visit(urls[topic] + strconv.Itoa(i))
	}
}
