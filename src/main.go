package main

import (
    "log"
    "net/http"
    cache "github.com/patrickmn/go-cache"
	"time"
)

var gocache = cache.New(5*time.Minute, 10*time.Minute)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/login", login)
    mux.HandleFunc("/logout", logout)
    mux.HandleFunc("/auth", auth)

    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}