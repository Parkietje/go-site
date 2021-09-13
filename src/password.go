package main

import (
    "golang.org/x/crypto/sha3"
	"encoding/hex"
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "errors"
)

const (
    hashesJson = "./data/hashes.json"
    saltsJson =  "./data/salts.json"
)

func passwordCheck(account string, password string) error {
    hashes, err := unmarshall(hashesJson)
    if err != nil{
        return errors.New("unauthorized")
    }

    salts, err := unmarshall(saltsJson)
    if err != nil{
        return errors.New("unauthorized")
    }

    var salt string

    if keyExists(salts, account){
        salt = salts[account].(string)
    }

    h := hash(password, salt)

    if hashes[account] != h{
        return errors.New("unauthorized")
    }

    return nil
}

func hash(password string, salt string) string{
    h := sha3.New512()
    h.Write([]byte(password+salt))
    sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func unmarshall(jsonFile string) (map[string]interface{}, error) {
    //read json file
    file, err := os.Open(jsonFile)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    defer file.Close()
    byteValue, _ := ioutil.ReadAll(file)

    // unmarshall the data
    var data map[string]interface{}
    err = json.Unmarshal(byteValue, &data)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return data, err
}

func keyExists(decoded map[string]interface{}, key string) bool {
    val, ok := decoded[key]
    return ok && val != nil
}