package website

import (
	"net/http"
	"text/template"

	"github.com/italo-carvalho/crowler/db"
)

type DataLinks struct {
	Links []db.VisitedLink
}

func indexHandle() func(http.ResponseWriter, *http.Request) {
	tmpl, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := db.FindAllLinks()
		if err != nil {
			panic(err)
		}
		data := DataLinks{Links: links}
		tmpl.Execute(w, data)
	}
}
