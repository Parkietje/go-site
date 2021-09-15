package main

import (
    "html/template"
    "log"           
    "net/http"
    "time"
    "errors"
    "github.com/google/uuid"
    "github.com/patrickmn/go-cache"
)

type Content struct {
    Navigation []Navitem
    PNG string
    User
}

type User struct {
    Account string
    SessionCookie string
}

type Navitem struct {
    Title string
    Route  string
}

const (
    HOME_template = "./ui/templates/home.gtpl"
    QR_template = "./ui/templates/qr.gtpl"
    LOGIN_template = "./ui/templates/login.gtpl"
)

var (
    DEFAULT_CONTENT = Content{Navigation: []Navitem{{Title: "Login", Route: "/login"}}}
    DEFAULT_USER = User{}
    DEFAULT_NAV = []Navitem{{Title: "Login", Route: "/login"}}
    AUTH_NAV = []Navitem{{Title: "Logout", Route: "/login"}}
    CACHE = cache.New(5*time.Minute, 10*time.Minute)
)

func home(w http.ResponseWriter, r *http.Request) {
    content := DEFAULT_CONTENT
    user, sc, err := verifySessionCookie(r)
    if err == nil {
        content.Navigation = AUTH_NAV
        img, _ := imgBase64Str("./ui/static/img/pngegg.png")
        content.PNG = img
        content.User = User{Account: user, SessionCookie: sc}
    }
    serve(HOME_template, content, w, r)
}

func login(w http.ResponseWriter, r *http.Request) {    
    if r.Method == "GET" {
        content := DEFAULT_CONTENT
        user, sc, err := verifySessionCookie(r)
        if err == nil {
            content.Navigation = AUTH_NAV
            content.User = User{Account: user, SessionCookie: sc}
        }
        serve(LOGIN_template, content, w, r)
    }

    if r.Method == "POST" {
        r.ParseForm()
        user := r.Form["username"][0]
        pw := r.Form["password"][0]
        
        if passwordCheck(user, pw) != nil{
            serve(LOGIN_template, DEFAULT_CONTENT, w, r)
        } else {
            secret := getSecret(user)
            // generate QR code
            img :=  genQR(user, secret)
            content := Content{DEFAULT_NAV, img, User{Account : user}}
            serve(QR_template, content, w, r)
        }
    }
}

func auth(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        token := r.Form["token"][0]
        account := r.Form["account"][0]
        secret := getSecret(account)

        _, err, authenticated := verify(token, secret)
        if authenticated && err == nil {
            // set session cookie for authenticated user
            token := uuid.New().String()
            setSessionCookie(account, token, w)
            img, _ := imgBase64Str("./ui/static/img/pngegg.png")
            content := Content{Navigation: AUTH_NAV , PNG: img, User: User{Account:account, SessionCookie: token}}
            serve(HOME_template, content, w, r)
            return
        }
    }
    home(w,r)
}

func logout(w http.ResponseWriter, r *http.Request){
    account, _, err := verifySessionCookie(r)
    if err == nil{
        deleteSessionCookie(account)
    }
    login(w,r)
}

func serve(file string, content Content, w http.ResponseWriter, r *http.Request) {
    files := []string{
        file,
        "./ui/templates/base.gtpl",
        "./ui/templates/footer.gtpl",
    }

    sc, account, err := verifySessionCookie(r)
    if err == nil {
        content.User = User{Account : account, SessionCookie: sc}
    }
    
    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }

    err = ts.Execute(w, content)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }
}

func setSessionCookie(account string, token string, w http.ResponseWriter){
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	CACHE.Set(account, token, 120 * time.Second)

	// set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
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

func deleteSessionCookie(account string){
	CACHE.Delete(account)
}

func verifySessionCookie(r *http.Request) (string, string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return "", "", errors.New("unauthorized")
	}
	sessionToken := c.Value

    c2, err := r.Cookie("account")
	if err != nil {
		return "", "", errors.New("unauthorized")
	}
	account := c2.Value

	st, _ := CACHE.Get(account)
	if st == nil {
		return account, "", errors.New("unauthorized")
	}

    if st != sessionToken {
        return account, "", errors.New("unauthorized")
    }

    return account, sessionToken, nil
}