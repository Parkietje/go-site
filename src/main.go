package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/term"
)

// Master secret to encrypt all generated qr codes.
// Do not store this in a file or env vars, instead remember it (!) or write it down.
// Master secret is only kept in memory for security purposes
// DO NOT FORGET THE MASTER SECRET OR YOU CANNOT RETRIEVE 2FA QR CODES
var MASTER_PASSWORD string

func main() {
	MASTER_PASSWORD = credentials()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/auth", auth)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// read credentials from stdin without echo'ing them in terminal history
// DO NOT LOG/PRINT SECRET
func credentials() string {
	fmt.Print("Enter master password: \n")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println(err)
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}
