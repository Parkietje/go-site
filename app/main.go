package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gosidekick/goconfig"
	"github.com/joho/godotenv"
)

type Config struct {
	ADMIN          string `cfgRequired:"true"`
	ADMIN_PASSWORD string `cfgRequired:"true"`
}

var (
	// used to check if user has admin rights
	ADMIN string
	// used as secret key for data encryption
	MASTER_PASSWORD string

	// holds our static web server content.
	//go:embed ui/static/*
	STATIC embed.FS
)

const (
	//Holds user authentication data
	HASHES = "data/hashes.json"
	//Holds user authentication data
	SALTS = "data/salts.json"
	//Holds user authentication data
	SECRETS = "data/secrets.json"
)

func main() {

	// load env vars to populate config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := Config{}
	err = goconfig.Parse(&config)
	if err != nil {
		println(err.Error())
		return
	}

	MASTER_PASSWORD = config.ADMIN_PASSWORD
	ADMIN = config.ADMIN

	// add admin credentials on first run
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777)
		var a [3]string
		a[0] = HASHES
		a[1] = SALTS
		a[2] = SECRETS
		for _, s := range a {
			f, err := os.Create(s)
			if err != nil {
				log.Fatal(err)
			}
			f.Write([]byte("{}"))
		}
		addUser(ADMIN, MASTER_PASSWORD, "some_s@lt")
	}

	r := mux.NewRouter()

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
