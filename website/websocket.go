package website

import (
	"context"
	"net/http"
	"text/template"
	"time"

	"github.com/italo-carvalho/crowler/crawler"
	"nhooyr.io/websocket"
)

func websocketHandle() func(http.ResponseWriter, *http.Request) {
	tmpl, err := template.ParseFiles("website/templates/websocket.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		website := r.FormValue("website")
		if website == "" {
			http.Error(w, "website field cannot be empty", http.StatusBadRequest)
			return
		}
		wc := crawler.New()
		go wc.VisitLink(website)

		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		subscriber(r.Context(), c, wc.Log())

		tmpl.Execute(w, nil)
	}
}

func subscriber(ctx context.Context, c *websocket.Conn, logs <-chan string) error {
	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-logs:
			err := writeTimeout(ctx, c, msg)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func writeTimeout(ctx context.Context, c *websocket.Conn, msg string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, []byte(msg))
}
