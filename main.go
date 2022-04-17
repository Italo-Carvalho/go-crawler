package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/italo-carvalho/crowler/db"
	"github.com/italo-carvalho/crowler/website"
	"golang.org/x/net/html"
)

var (
	link   string
	action string
)

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "init url crawler")
	flag.StringVar(&action, "action", "website", "what service to start?")
}

func main() {
	flag.Parse()
	switch *&action {
	case "website":
		website.Run()
	case "crawler":
		done := make(chan bool)
		go visitLink(link)
		<-done
	default:
		fmt.Printf("unknown action: '%s'", action)
	}
}

func visitLink(link string) {
	fmt.Printf("visitando: %s \n", link)
	resp, err := http.Get(link)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close() // fechar conexao

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ERROR] status diferente de 200: %d\n", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		panic(err)
	}

	extractLinks(doc)
}

func extractLinks(node *html.Node) {
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
				fmt.Printf("link is already visited: %s \n", link)
				continue
			}

			visitedLink := db.VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			db.Insert("links", visitedLink)

			go visitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}
