package main

import (
	"log"
	"net/http"
	"os"
	"./pixels"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/pixels", pixels.PixelsHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
