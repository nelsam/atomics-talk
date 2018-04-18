package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %s", err)
	}
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(wd))))
}
