package main

import (
	"flag"
	"log"
	"net/http"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

    r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
        c := r.URL.Query().Get("channel")
		serveWs(hub, w, r, c)
	})

	err := http.ListenAndServe(*addr, handlers.CORS()(r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
