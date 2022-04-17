package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/italo-carvalho/crowler/db"
	"golang.org/x/net/html"
)

type webCrawler struct {
	log chan string
}

func (wc *webCrawler) VisitLink(link string) {
	wc.log <- fmt.Sprintf("visitando: %s", link)
	resp, err := http.Get(link)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close() // fechar conexao

	if resp.StatusCode != http.StatusOK {
		wc.log <- fmt.Sprintf("[ERROR] The status is defferent from 200: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		panic(err)
	}

	wc.extractLinks(doc)
}

func (wc *webCrawler) extractLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}
			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" || !strings.HasPrefix(link.String(), "http") {
				continue
			}

			if db.CheckVisitedLink(link.String()) {
				wc.log <- fmt.Sprintf("link is already visited: %s", link)
				continue
			}

			visitedLink := db.VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			db.Insert("links", visitedLink)

			go wc.VisitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		wc.extractLinks(c)
	}
}

func (wc *webCrawler) Log() chan string {
	return wc.log
}

func New() *webCrawler {
	wb := &webCrawler{
		log: make(chan string, 10),
	}

	return wb
}
