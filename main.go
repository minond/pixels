package main

import (
	"log"
	"net/http"
	"os"
	"github.com/minond/pixels"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/pixels", PixelsHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
