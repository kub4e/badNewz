package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getContent(url string, filter string) string {
	result := ""
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	checkErr(err)

	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36")

	resp, err := client.Do(request)
	checkErr(err)
	defer resp.Body.Close()

	document, _ := goquery.NewDocumentFromReader(resp.Body)
	document.Find(filter).Each(func(i int, s *goquery.Selection) {

		paragraphs := make([]string, 0)
		s.Find("p").Each(func(j int, k *goquery.Selection) {
			paragraphs = append(paragraphs, k.Text())
		})
		result = strings.Join(paragraphs, "\n")
	})
	return result
}
