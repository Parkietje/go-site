package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type Context struct {
	Navigation []Navitem
	PNG        string
	User
}

type User struct {
	Account       string
	SessionCookie string
}

type Navitem struct {
	Title string
	Route string
}

const (
	HOME_template  = "./ui/templates/home.gtpl"
	QR_template    = "./ui/templates/qr.gtpl"
	LOGIN_template = "./ui/templates/login.gtpl"
)

var (
	DEFAULT_CONTEXT = Context{Navigation: []Navitem{{Title: "Login", Route: "/login"}}}
	AUTH_NAV        = []Navitem{{Title: "Logout", Route: "/login"}}
	CACHE           = cache.New(5*time.Minute, 10*time.Minute)
)

func home(w http.ResponseWriter, r *http.Request) {
	context, auth := getContext(r)
	if auth {
		img, _ := imgBase64Str("./ui/static/img/pngegg.png")
		context.PNG = img
	}
	render(HOME_template, context, w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context, _ := getContext(r)
		render(LOGIN_template, context, w, r)
	}
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form["username"][0]
		pw := r.Form["password"][0]
		// generate QR code if password matches
		if passwordCheck(user, pw) != nil {
			render(LOGIN_template, DEFAULT_CONTEXT, w, r)
		} else {
			secret := getSecret(user)
			img := genQR(user, secret)
			context := Context{AUTH_NAV, img, User{Account: user}}
			render(QR_template, context, w, r)
		}
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		token := r.Form["token"][0]
		account := r.Form["account"][0]
		secret := getSecret(account)
		// set session cookie for authenticated user
		_, err, authenticated := verify(token, secret)
		if authenticated && err == nil {
			token := uuid.New().String()
			setSessionCookie(account, token, w)
			img, _ := imgBase64Str("./ui/static/img/pngegg.png")
			context := Context{Navigation: AUTH_NAV, PNG: img, User: User{Account: account, SessionCookie: token}}
			render(HOME_template, context, w, r)
			return
		}
	}
	home(w, r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	account, _, err := verifySessionCookie(r)
	if err == nil {
		deleteSessionCookie(account, w)
	}
	login(w, r)
}

func getContext(r *http.Request) (Context, bool) {
	context := DEFAULT_CONTEXT
	sc, account, err := verifySessionCookie(r)
	if err != nil {
		// user not authenticated, return default context
		return context, false
	} else {
		// user authenticated, return authenticated context
		context.Navigation = AUTH_NAV
		context.User = User{Account: account, SessionCookie: sc}
		return context, true
	}
}

func render(file string, context Context, w http.ResponseWriter, r *http.Request) {
	files := []string{
		file,
		"./ui/templates/base.gtpl",
		"./ui/templates/footer.gtpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, context)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func setSessionCookie(account string, token string, w http.ResponseWriter) {
	// add username to cache
	if account != "" {
		CACHE.Set(account, token, 120*time.Second)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(120 * time.Second),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "account",
		Value:   account,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func deleteSessionCookie(account string, w http.ResponseWriter) {
	//remove token from cache
	CACHE.Delete(account)
	//delete request cookie by setting empty value
	setSessionCookie("", "", w)
}

func verifySessionCookie(r *http.Request) (string, string, error) {
	c, err := r.Cookie("account")
	if err != nil {
		return "", "", errors.New("unauthorized")
	}
	c2, err := r.Cookie("session_token")
	if err != nil {
		return "", "", errors.New("unauthorized")
	}

	//check if request cookie is stored in cache
	account := c.Value
	cookie_token := c2.Value
	cache_token, _ := CACHE.Get(account)

	if cache_token == nil || cache_token == "" {
		return account, "", errors.New("unauthorized")
	}
	if cache_token != cookie_token {
		return account, "", errors.New("unauthorized")
	}
	return account, cookie_token, nil
}
