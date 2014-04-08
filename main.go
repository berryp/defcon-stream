package main

import (
	"fmt"
	"log"
	"net/http"
)

const staticRoot = "app"

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) error {
	var path = r.URL.Path

	if path == "/" {
		path = "/index.html"
	}

	http.ServeFile(w, r, staticRoot+path)

	return nil
}

func main() {
	http.Handle("/", appHandler(mainHandler))

	fmt.Println("Listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
