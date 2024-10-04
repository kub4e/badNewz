package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

func generateSites(articleList []article) {
	for _, article := range articleList {
		t, err := template.ParseFiles("../shared/template.html")
		if err != nil {
			log.Print(err)
			return
		}
		name := article.Url
		name = trimUrl(name)

		f, err := os.Create("../shared/articles/" + name)
		if err != nil {
			log.Println("create file: ", err)
			return
		}
		defer f.Close()

		err = t.Execute(f, article)
		if err != nil {
			log.Print("execute: ", err)
			return
		}
	}
}

func trimUrl(name string) string {

	if name[len(name)-1:] == "/" {
		name = name[:len(name)-1]
	}

	if name[len(name)-4:] == "html" {
		name = strings.Split(name, "/")[len(strings.Split(name, "/"))-1]
	} else {
		name = strings.Split(name, "/")[len(strings.Split(name, "/"))-1] + ".html"
	}
	return name
}
