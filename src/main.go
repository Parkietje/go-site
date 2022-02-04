package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Master secret for data encryption
var MASTER_PASSWORD string

// The admin account name
var ADMIN = "bd6e0ef5756b899a944956e023e3d8a10504955b9bb3dc831287d9985723c5428a180e75ddb2286379302e02e948faefd032aa8b492187e0a9e7e7ee0420d6e3"

// holds our static web server content.
//go:embed ui/static/*
var STATIC embed.FS

func main() {
	adminFlag := flag.Bool("admin", false, "if true an admin account can be added on startup")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MASTER_PASSWORD = os.Getenv("MASTER_PASS")

	if *adminFlag {
		addAdmin()
	}

	r := mux.NewRouter()

	// TODO: embed artifacts in binary
	// https://stackoverflow.com/questions/13904441/whats-the-best-way-to-bundle-static-resources-in-a-go-program
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(STATIC))))
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/admin", admin)
	r.HandleFunc("/admin/{service}", admin)

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", r)
	log.Fatal(err)
}
