package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/morserino", morserino)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
