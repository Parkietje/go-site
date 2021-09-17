package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Messages   []Message
	Webstats
	PNG string
}

type Webstats struct {
	Stats []Stat
	IPs   []IP
}

type IP struct {
	Address string
	Count   string
}

type Stat struct {
	Name  string
	Value string
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
	HOME_template  = "../ui/templates/home.gtpl"
	QR_template    = "../ui/templates/qr.gtpl"
	LOGIN_template = "../ui/templates/login.gtpl"
	ADMIN_template = "../ui/templates/admin.gtpl"
)

var (
	DEFAULT_NAV     = []Navitem{{Title: "Login", Route: "/login"}}
	AUTH_NAV        = []Navitem{{Title: "Logout", Route: "/login"}}
	ADMIN_NAV       = []Navitem{{Title: "Admin", Route: "/admin"}, {Title: "Logout", Route: "/login"}}
	DEFAULT_CONTENT = PageContent{Navigation: DEFAULT_NAV}
	DEFAULT_CONTEXT = Context{User{}, DEFAULT_CONTENT}
	CACHE           = cache.New(5*time.Minute, 10*time.Minute)
)

func home(w http.ResponseWriter, r *http.Request) {
	addStat("page_visits")
	addIP(ReadUserIP(r))
	context, auth := getContext(r)
	if auth {
		img, _ := imgBase64Str("../ui/static/img/pngegg.png")
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
		if passwordCheck(hash(user, ""), pw) != nil {
			addStat("failed_logins")
			render(LOGIN_template, DEFAULT_CONTEXT, w, r)
		} else {
			addStat("logins")
			secret := getSecret(hash(user, ""))
			img := genQR(user, secret)
			context := Context{User{Account: hash(user, "")}, PageContent{AUTH_NAV, nil, Webstats{}, img}}
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
			img, _ := imgBase64Str("../ui/static/img/pngegg.png")
			context := Context{User{account, token}, PageContent{AUTH_NAV, nil, Webstats{}, img}}
			if account == ADMIN {
				context.PageContent.Navigation = ADMIN_NAV
			}
			render(HOME_template, context, w, r)
			return
		}
	}
	home(w, r)
}

func admin(w http.ResponseWriter, r *http.Request) {
	user, _, err := verifySessionCookie(r)
	if (err == nil) && (user == ADMIN) {
		context, _ := getContext(r)
		urlparts := strings.Split(r.RequestURI, "/")
		switch service := urlparts[2]; service {

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

		default:
			fmt.Println("default")
		}

		ips, _ := unmarshal(IP_FILE)
		var ipslice []IP
		for k, v := range ips {
			var ip IP
			ip.Address = k
			ip.Count = v
			ipslice = append(ipslice, ip)
		}
		stats, _ := unmarshal(STATS_FILE)
		wstats := []Stat{{"logins", stats["logins"]}, {"failed_logins", stats["failed_logins"]}, {"page_visits", stats["page_visits"]}}
		context.PageContent.Webstats.Stats = wstats
		context.PageContent.Webstats.IPs = ipslice

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

func render(file string, context Context, w http.ResponseWriter, r *http.Request) {
	files := []string{
		file,
		"../ui/templates/base.gtpl",
		"../ui/templates/footer.gtpl",
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

func deleteSessionCookie(account string, w http.ResponseWriter) {
	//remove token from cache
	CACHE.Delete(account)
	//delete request cookie by setting empty value
	setSessionCookie("", "", w)
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
