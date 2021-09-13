package main

import (
    "html/template"
    "log"           
    "net/http"
    guuid "github.com/google/uuid"
    "time"
    "errors"
    "encoding/json"
    "os"
    "io/ioutil"
    "fmt"
)

const (
    authorizedUsers = "./data/accounts.json"
)

type Navitem struct {
    Title string
    Route  string
}

type Content struct {
    Nav []Navitem
    Img string
    Account string
    SessionCookie string
}

func home(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/templates/home.gtpl",
        "./ui/templates/base.gtpl",
        "./ui/templates/footer.gtpl",
    }
    
    _, _, err := verifySessionCookie(r)
    
    if err != nil{
        serve(files, w, r, Content{Nav: nil, Img: ""})
    } else {
        data := []Navitem{{Title: "authenticated", Route: "/auth"}}
        img, err := imgBase64Str(penguinFilename)
        if err != nil {
            img = ""
        }
        content := Content{Nav : data, Img: img}
        serve(files, w, r, content)
    }
}

func login(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/templates/login.gtpl",
        "./ui/templates/base.gtpl",
        "./ui/templates/footer.gtpl",
    }
    
    if r.Method == "GET" {
        sc, account, err := verifySessionCookie(r)
        if err != nil{
            serve(files, w, r, Content{Nav: nil, Img: ""})
        } else {
            data := []Navitem{{Title: "authenticated", Route: "/auth"}}
            content := Content{Nav : data, Account : account, SessionCookie: sc}
            serve(files, w, r, content)
        }
    }

    if r.Method == "POST" {
        r.ParseForm()
        user := r.Form["username"][0]
        pw := r.Form["password"][0]
        if passwordCheck(user, pw) != nil{
            serve(files, w, r, Content{Nav: nil, Img: ""})
        } else {
            // generate QR code
            base64str :=  genQR(user)
            content := Content{Nav : nil, Img: base64str, Account: user}
            files := []string{
                "./ui/templates/qr.gtpl",
                "./ui/templates/base.gtpl",
                "./ui/templates/footer.gtpl",
            }
            serve(files, w, r, content)
        }
    }
}

func passwordCheck(account string, password string) error {
    jsonFile, err := os.Open(authorizedUsers)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    // unmarshall the data
    var data map[string]interface{}
    err = json.Unmarshal(byteValue, &data)
    if err != nil {
        fmt.Println(err)
        return err
    }

    h := hash(password)

    if data[account] != h{
        return errors.New("unauthorized")
    }

    return nil
}

func logout(w http.ResponseWriter, r *http.Request){
    account, _, err := verifySessionCookie(r)
    if err == nil{
        deleteSessionCookie(account)
    }
    login(w,r)
}

func auth(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        token := r.Form["token"]
        tokenstr := token[0]
        account := r.Form["account"]
        accountstr := account[0] 
       
        _, err, authenticated := verify(tokenstr)
        if authenticated && err == nil {
            files := []string{
                "./ui/templates/authenticated.gtpl",
                "./ui/templates/base.gtpl",
                "./ui/templates/footer.gtpl",
            }

            // set session cookie for authenticated user
            token := guuid.New().String()
            setSessionCookie(accountstr, token, w)

            data := []Navitem{{Title: "authenticated", Route: "/auth"}}
            content := Content{Nav : data}
            serve(files, w, r, content)
        } else {
            fmt.Println(err)
            home(w,r)
        }
    }

    if r.Method == "GET" {
        _, _, err := verifySessionCookie(r)
        if err != nil{
            home(w,r)
        } else {
            files := []string{
                "./ui/templates/authenticated.gtpl",
                "./ui/templates/base.gtpl",
                "./ui/templates/footer.gtpl",
            }
            data := []Navitem{{Title: "authenticated", Route: "/auth"}}
            content := Content{Nav : data}
            serve(files, w, r, content)
        }
    }
}

func serve(files []string, w http.ResponseWriter, r *http.Request, data Content) {
    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
        return
    }

    err = ts.Execute(w, data)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }
}

func setSessionCookie(account string, token string, w http.ResponseWriter){
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	gocache.Set(account, token, 120 * time.Second)

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
	gocache.Delete(account)
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

	st, _ := gocache.Get(account)
	if st == nil {
		return account, "", errors.New("unauthorized")
	}

    if st != sessionToken {
        return account, "", errors.New("unauthorized")
    }

    return account, sessionToken, nil
}