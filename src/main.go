package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Master secret to encrypt all generated qr codes.
// Do not store this in a file or env vars, instead remember it (!) or write it down.
// Master secret is only kept in memory for security purposes
// DO NOT FORGET THE MASTER SECRET OR YOU CANNOT RETRIEVE 2FA QR CODES
var MASTER_PASSWORD string

func main() {
	adminFlag := flag.Bool("admin", false, "if true an admin account can be added on startup")
	flag.Parse()

	MASTER_PASSWORD = readPassword()

	if *adminFlag {
		addAdmin()
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../ui/static/"))))
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/auth", auth)
	r.HandleFunc("/admin/", admin)
	r.HandleFunc("/admin/{service}", admin)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", r)
	log.Fatal(err)
}
