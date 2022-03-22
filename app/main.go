package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gosidekick/goconfig"
	"github.com/joho/godotenv"
)

type Config struct {
	ADMIN           string `cfgRequired:"true"`
	ADMIN_PASSWORD  string `cfgRequired:"true"`
	STORAGE_ACCOUNT string `cfgRequired:"False"`
	CONTAINER_NAME  string `cfgRequired:"False"`
	STORAGE_KEY     string `cfgRequired:"False"`
}

var (
	// used to check if user has admin rights
	ADMIN string
	// used as secret key for data encryption
	MASTER_PASSWORD string
	// holds our static web server content.
	//go:embed ui/*
	FS        embed.FS
	static, _ = fs.Sub(FS, "ui/static")
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
	STORAGE_ACCOUNT = config.STORAGE_ACCOUNT
	STORAGE_KEY = config.STORAGE_KEY
	CONTAINER_NAME = config.CONTAINER_NAME

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

	// UI ROUTES
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	r.HandleFunc("/", home)
	r.HandleFunc("/admin", admin)
	r.HandleFunc("/admin/{service}", admin)
	r.HandleFunc("/deploy", deploy)
	r.HandleFunc("/deploy/{command}", deploy)
	r.HandleFunc("/logout", logout)

	// API ROUTES
	r.HandleFunc("/auth", auth)
	r.HandleFunc("/list", list)
	r.HandleFunc("/upload", upload)

	log.Println("Starting server on :8000")
	err = http.ListenAndServe(":8000", r)
	log.Fatal(err)
}
