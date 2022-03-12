package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func list(w http.ResponseWriter, r *http.Request) {
	_, auth := getContext(r)
	if auth {
		s, e := listBlobs()
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

func upload(w http.ResponseWriter, r *http.Request) {
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
		uploadBlob(path)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}