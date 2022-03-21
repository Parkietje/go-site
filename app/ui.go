package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type Context struct {
	User
	PageContent
}

type User struct {
	Account       string
	SessionCookie string
}

type PageContent struct {
	Navigation []Navitem
	Sidebar    []Navitem
	PNG        string
}

type Navitem struct {
	Title string
	Route string
}

const (
	HOME_template   = "ui/pages/home.gohtml"
	ADMIN_template  = "ui/pages/admin.gohtml"
	DEPLOY_template = "ui/pages/deploy.gohtml"
)

var (
	DEFAULT_NAV     = []Navitem{}
	AUTH_NAV        = []Navitem{{Title: "Deployments", Route: "/deploy"}}
	ADMIN_NAV       = []Navitem{{Title: "Deployments", Route: "/deploy"}, {Title: "Admin", Route: "/admin"}}
	DEFAULT_CONTENT = PageContent{Navigation: DEFAULT_NAV}
	DEFAULT_CONTEXT = Context{User{}, DEFAULT_CONTENT}
)

func home(w http.ResponseWriter, r *http.Request) {
	context, auth := getContext(r)
	if auth {
		context.Sidebar = []Navitem{{Title: "Log out", Route: "login"}, {Title: "Image", Route: "img"}}
		context.PageContent.PNG, _ = imgBase64Str("ui/static/icons/pngegg.png")
	} else {
		context.Sidebar = []Navitem{{Title: "Login", Route: "login"}}
	}
	render(HOME_template, context, w, r)
}

func admin(w http.ResponseWriter, r *http.Request) {
	user, _, err := verifySessionCookie(r)
	if (err == nil) && (user == ADMIN) {
		context, _ := getContext(r)
		urlparts := strings.Split(r.RequestURI, "/")
		var service string
		if len(urlparts) >= 3 {
			service = urlparts[2]
		}
		switch service {

		case "add":
			r.ParseForm()
			user := r.Form["username"][0]
			password := r.Form["password"][0]
			salt := hash(genSecret(), "s@lty?")
			addUser(user, password, salt)
			fmt.Println("user added")

		case "delete":
			r.ParseForm()
			user := r.Form["hash"][0]
			deleteUser(user)
			fmt.Println("user deleted")

		}

		render(ADMIN_template, context, w, r)
	} else {
		home(w, r)
	}
}

func deploy(w http.ResponseWriter, r *http.Request) {
	context, auth := getContext(r)
	if auth {
		render(DEPLOY_template, context, w, r)
	} else {
		render(HOME_template, context, w, r)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	account, _, err := verifySessionCookie(r)
	if err == nil {
		deleteSessionCookie(account, w)
		fmt.Println("session cookie deleted")
	}
	http.Redirect(w, r, "http://localhost:8000", http.StatusSeeOther)
}

// parse html templates and execute response
func render(file string, context Context, w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(FS, file, "ui/layout/*.gohtml", "ui/components/*.gohtml")
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

// get page context for current user
func getContext(r *http.Request) (Context, bool) {
	context := DEFAULT_CONTEXT
	account, sc, err := verifySessionCookie(r)
	if err != nil {
		return context, false
	} else {
		context.User = User{Account: account, SessionCookie: sc}
		if account == ADMIN {
			context.PageContent.Navigation = ADMIN_NAV
			return context, true
		} else {
			context.PageContent.Navigation = AUTH_NAV
			return context, true
		}
	}
}

//encode PNG to html-embeddable string
func imgBase64Str(fileName string) (string, error) {
	f, err := FS.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(f), nil
}

// cache to hold session cookies
var CACHE = cache.New(5*time.Minute, 15*time.Minute)

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

func setSessionCookie(account string, token string, w http.ResponseWriter) {
	// add username to cache
	if account != "" {
		CACHE.Set(account, token, 15*time.Minute)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(15 * time.Minute),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "account",
		Value:   account,
		Expires: time.Now().Add(15 * time.Minute),
	})
}

func deleteSessionCookie(account string, w http.ResponseWriter) {
	//remove token from cache
	CACHE.Delete(account)
	//delete request cookie by setting empty value
	setSessionCookie("", "", w)
}
