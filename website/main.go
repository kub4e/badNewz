package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type article struct {
	Title       string
	Description string
	Date        string
	Url         string
	Img         string
	Topic       string
	Source      string
	Content     string
	MyUrl       string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	r.PathPrefix("/articles/").Handler(http.StripPrefix("/articles/", http.FileServer(http.Dir("../shared/articles"))))
	http.ListenAndServe(":3000", r)
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	//articleList = make([]string, 0)

	topic := req.URL.Query().Get("topic")
	source := req.URL.Query().Get("source")
	if topic == "" || topic == "all" {
		topic = "%%"
	}
	if source == "" || source == "all" {
		source = "%%"
	}

	page := req.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)
	articlesPerPage := 20
	firstIndex := (pageInt - 1) * articlesPerPage

	search := req.URL.Query().Get("search")

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, loadFromDB(topic, source, search, firstIndex))
	if err != nil {
		panic(err)
	}
}

func loadFromDB(topicFilter, sourceFilter, searchFilter string, offset int) []article {
	articleList := make([]article, 0)
	db, err12 := sql.Open("sqlite3", "../shared/articles.db")
	checkErr(err12)
	rows, _ := db.Query("SELECT * FROM articles WHERE topic like ? AND source like ? limit 20 offset ?", topicFilter, sourceFilter, offset)

	var title, description, date, url, img, topic, source, content string
	for rows.Next() {
		rows.Scan(&title, &description, &date, &url, &img, &topic, &source, &content)

		if !strings.Contains(title, searchFilter) {
			continue
		}

		a := article{
			Title:       title,
			Description: description,
			Date:        date,
			Url:         url,
			Img:         img,
			Topic:       topic,
			Source:      source,
			Content:     content,
		}
		a.MyUrl = trimUrl(a.Url)
		articleList = append(articleList, a)
	}
	rows.Close()
	return articleList
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
