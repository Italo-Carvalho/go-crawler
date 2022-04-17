package website

import (
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", indexHandle())
	http.HandleFunc("/busca", websocketHandle())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
