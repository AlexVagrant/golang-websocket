package main

import (
	_ "fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader:", err)
		return
	}
	defer conn.Close()
	for {
		// receive message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("recv: %s", p)
		// send message
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {

	serverMux := http.NewServeMux()

	s := &http.Server{
		Addr:         ":7070",
		Handler:      serverMux,
		WriteTimeout: 10 * time.Second,
	}

	serverMux.HandleFunc("/", home)
	serverMux.HandleFunc("/ws", ws)
	serverMux.Handle("/static", http.FileServer(http.Dir("/static")))

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
