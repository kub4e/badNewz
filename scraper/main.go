// Make the functions get a single article

package main

import (
	"database/sql"
	"sync"

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
}

func main() {
	articleList := make([]article, 0)
	depth := 2

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		scrapeNoviniBG(depth, &articleList, "politics")
		scrapeNoviniBG(depth, &articleList, "economy")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scrapeBNT(depth, &articleList, "economy")
		scrapeBNT(depth, &articleList, "politics")
		scrapeBNT(depth, &articleList, "society")
		scrapeBNT(depth, &articleList, "justice")
		wg.Done()
	}()

	/*
		Website no longer maintained.
		wg.Add(1)
		go func() {
			scrapeTVEVROPA(depth, &articleList, "politics")
			scrapeTVEVROPA(depth, &articleList, "economy")
			wg.Done()
		}()
	*/

	wg.Add(1)
	go func() {
		scrapteBTV(depth, &articleList)
		wg.Done()
	}()

	wg.Wait()
	//fmt.Println(articleList)

	for i, j := 0, len(articleList)-1; i < j; i, j = i+1, j-1 {
		articleList[i], articleList[j] = articleList[j], articleList[i]
	}

	db, err := sql.Open("sqlite3", "../shared/articles.db")
	checkErr(err)
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS articles (title TEXT PRIMARY KEY,
		description TEXT,
		date TEXT,
		url TEXT,
		img TEXT,
		topic TEXT,
		source TEXT,
		content TEXT)`)
	checkErr(err)
	statement.Exec()

	statement, err = db.Prepare("INSERT INTO articles (title, description, date, url, img, topic, source, content) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)

	for _, article := range articleList {
		//fmt.Println(reflect.TypeOf(article.Content))
		statement.Exec(article.Title, article.Description, article.Date, article.Url, article.Img, article.Topic, article.Source, article.Content)
	}

	generateSites(articleList)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
