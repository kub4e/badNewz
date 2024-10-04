package main

import (
	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

func scrapeBNT(depth int, articleList *[]article, topic string) {
	urls := make(map[string]string)
	urls["politics"] = "https://bntnews.bg/bg/c/bgpolitika?page="
	urls["economy"] = "https://bntnews.bg/bg/c/bgikonomika?page="
	urls["society"] = "https://bntnews.bg/bg/c/obshtestvo-15?page="
	urls["justice"] = "https://bntnews.bg/bg/c/bgsigurnost?page="

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"),
	)

	c.OnHTML(".news-wrap-view > a[href], .img-wrap > a[href]", func(e *colly.HTMLElement) {

		a := article{
			Title:  e.Attr("title"),
			Url:    e.Attr("href"),
			Date:   e.ChildText(".news-time"),
			Img:    e.ChildAttr("img", "src"),
			Topic:  topic,
			Source: "BNT",
		}
		a.Content = getContent(a.Url, ".txt-news")
		*articleList = append(*articleList, a)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i <= depth; i++ {
		c.Visit(urls[topic] + strconv.Itoa(i))
	}
}
