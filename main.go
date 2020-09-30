package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(page string) {
	hostname := "https://studygolang.com"
	// Request the HTML page.
	res, err := http.Get(hostname + "/go/weekly?p=" + page)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if l, ok := s.Attr("href"); ok &&
			strings.HasPrefix(l, "/topics/") &&
			!strings.Contains(l, "#") &&
			!strings.HasPrefix(l, "/topics/new") {

			// Request the HTML page.
			res, err := http.Get(hostname + l)
			if err != nil {
				log.Fatal(err)
			}
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}

			// Load the HTML document
			topic, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			res.Body.Close()

			findWhat := false
			topic.Find("div.content.markdown-body").Each(func(j int, content *goquery.Selection) {
				findWhat = strings.Contains(content.Text(), "what")
			})
			if findWhat {
				fmt.Println("找到what!!! ", s.Text(), l)
			} else {
				fmt.Println(s.Text(), l)
			}
		}
	})
}

func main() {
	ExampleScrape("1")
	ExampleScrape("2")
}
