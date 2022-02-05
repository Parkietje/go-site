package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

// cache to hold session cookies
var CACHE = cache.New(5*time.Minute, 10*time.Minute)

func home(w http.ResponseWriter, r *http.Request) {
	context, auth := getContext(r)
	if auth {
		context.PageContent.PNG, _ = imgBase64Str("ui/static/img/pngegg.png")
	}
	render(HOME_template, context, w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context, auth := getContext(r)
		if auth {
			context.PageContent.PNG, _ = imgBase64Str("ui/static/img/pngegg.png")
		}
		render(LOGIN_template, context, w, r)
	}
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form["username"][0]
		pw := r.Form["password"][0]
		if passwordCheck(hash(user, ""), pw) != nil {
			render(LOGIN_template, DEFAULT_CONTEXT, w, r)
		} else {
			token := uuid.New().String()
			setSessionCookie(user, token, w)
			img, _ := imgBase64Str("ui/static/img/pngegg.png")
			context := Context{User{user, token}, PageContent{AUTH_NAV, nil, img}}
			if user == ADMIN {
				context.PageContent.Navigation = ADMIN_NAV
			}
			render(HOME_template, context, w, r)
		}
	}
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
			salt := r.Form["salt"][0]
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

func logout(w http.ResponseWriter, r *http.Request) {
	account, _, err := verifySessionCookie(r)
	if err == nil {
		deleteSessionCookie(account, w)
	}
	login(w, r)
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