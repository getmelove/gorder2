package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	log.Println("Listening on :8082")
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.RequestURI)
		io.WriteString(w, "pong")
	})
	if err := http.ListenAndServe(":8082", mux); err != nil {
		panic(err)
	}
}
