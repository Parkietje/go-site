package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
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

type Message struct {
	Title string
	Text  string
}

const (
	HOME_template   = "ui/pages/home.gohtml"
	LOGIN_template  = "ui/pages/login.gohtml"
	ADMIN_template  = "ui/pages/admin.gohtml"
	DEPLOY_template = "ui/pages/deploy.gohtml"
)

var (
	DEFAULT_NAV     = []Navitem{{Title: "Home", Route: "/"}, {Title: "Login", Route: "/login"}}
	AUTH_NAV        = []Navitem{{Title: "Home", Route: "/"}, {Title: "Deployments", Route: "/deploy"}}
	ADMIN_NAV       = []Navitem{{Title: "Home", Route: "/"}, {Title: "Deployments", Route: "/deploy"}, {Title: "Admin", Route: "/admin"}}
	DEFAULT_CONTENT = PageContent{Navigation: DEFAULT_NAV}
	DEFAULT_CONTEXT = Context{User{}, DEFAULT_CONTENT}
)

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

//encode PNG to html-embeddable string
func imgBase64Str(fileName string) (string, error) {
	f, err := FS.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(f), nil
}
