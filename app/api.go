package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type loginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
		var u loginUser
		if err := json.Unmarshal(body, &u); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
		user := u.Username
		pw := u.Password
		if passwordCheck(hash(user, ""), pw) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		} else {
			token := uuid.New().String()
			setSessionCookie(user, token, w)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(token))
		}
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	user, _, err := verifySessionCookie(r)
	if (err == nil) && (user == ADMIN) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
		var u loginUser
		if err := json.Unmarshal(body, &u); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}

		switch r.Method {
		case http.MethodPost:
			salt := hash(genSecret(), "s@lty?")
			err = saveUser(u.Username, u.Password, salt)
			if err == nil {
				fmt.Println("user added")
			} else {
				fmt.Println("couldn't add user")
			}
			return

		case http.MethodDelete:
			err = removeUser(u.Username)
			if err == nil {
				fmt.Println("user deleted")
			} else {
				fmt.Println("couldn't delete user")
			}
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}

func listBlobs(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		s, e := listAzureBlobs()
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		w.Write([]byte(s))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}

func uploadBlob(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		r.ParseMultipartForm(32 << 20) // limit your max input length!
		var buf bytes.Buffer
		file, header, err := r.FormFile("files[]")
		if err != nil {
			fmt.Printf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}
		defer file.Close()
		name := header.Filename
		fmt.Printf("File name %s\n", name)
		io.Copy(&buf, file)
		contents := buf.String()
		if _, err := os.Stat("data/blobs"); os.IsNotExist(err) {
			os.Mkdir("data/blobs", 0777)
		}
		path := filepath.Join("data/blobs", name)
		ioutil.WriteFile(path, []byte(contents), 0644)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		buf.Reset()
		uploadAzureBlob(path)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}

func listMongos(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		s, e := listVMs()
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		w.Write([]byte(s))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}

func createMongo(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		mongos, _ := listVMs()
		var arr []string
		_ = json.Unmarshal([]byte(mongos), &arr)
		count := len(arr)
		name := "mongo-test-" + fmt.Sprintf("%d", count)
		fmt.Println(name)
		bindIP, _ := deployAzureMongo(name)
		fmt.Println(bindIP)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}

type mongoSetupRequest struct {
	InstanceName  string `json:"name"`
	DBName        string `json:"dbname"`
	MongoUser     string `json:"user"`
	MongoPassword string `json:"pass"`
}

func setupMongo(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}
		var msr mongoSetupRequest
		if err := json.Unmarshal(body, &msr); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}

		fmt.Println(msr.InstanceName)
		IP, _ := getIP(msr.InstanceName)
		fmt.Println(IP)

		SCP("C:/Users/YanniChiodi/Code/go-site/mongo-setup/", IP)

		execute("chmod +x /home/ubuntu/install_mongodb.sh", IP)
		execute("sudo /home/ubuntu/install_mongodb.sh", IP)

		execute("chmod +x /home/ubuntu/setup_user_and_db_mongo.sh", IP)
		execute("sudo /home/ubuntu/setup_user_and_db_mongo.sh -u="+msr.MongoUser+" -p="+msr.MongoPassword+" -d="+msr.DBName, IP)

		execute("chmod +x /home/ubuntu/set_mongo_bindIp.sh", IP)
		execute("sudo /home/ubuntu/set_mongo_bindIp.sh", IP)

		execute("chmod +x /home/ubuntu/ping.sh", IP)
		_, e := execute("sudo /home/ubuntu/ping.sh", IP)
		if e == nil {
			fmt.Println("mongo setup successful")
		} else {
			fmt.Println("an error occured while setting up the mongo")
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}

func checkMongoStatus(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		urlparts := strings.Split(r.RequestURI, "/")
		instance := urlparts[3]
		IP, _ := getIP(instance)
		_, e := execute("sudo /home/ubuntu/ping.sh", IP)
		w.WriteHeader(http.StatusOK)
		if e == nil {
			fmt.Println("mongo running")
			w.Write([]byte("mongo running"))
		} else {
			fmt.Println("mongo down")
			w.Write([]byte("mongo down"))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}
