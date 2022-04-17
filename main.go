package main

import (
	"flag"
	"fmt"

	"github.com/italo-carvalho/crowler/crawler"
	"github.com/italo-carvalho/crowler/website"
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
		wb := crawler.New()
		go wb.VisitLink(link)
		for log := range wb.Log() {
			fmt.Println(log)
		}
		<-done
	default:
		fmt.Printf("unknown action: '%s'", action)
	}
}
