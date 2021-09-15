package main

import (
    "golang.org/x/crypto/sha3"
    "encoding/hex"
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "errors"
    "crypto/rand"
    "encoding/base32"
)

const (
    JSON_hashes = "./data/hashes.json"
    JSON_salts =  "./data/salts.json"
    JSON_secrets = "./data/secrets.json"
)

func hash(password string, salt string) string{
    h := sha3.New512()
    h.Write([]byte(password+salt))
    sum := h.Sum(nil)
    return hex.EncodeToString(sum)
}

func passwordCheck(account string, password string) error {
    hashes, err := unmarshal(JSON_hashes)
    if err != nil{
        return errors.New("unauthorized")
    }

    salts, err := unmarshal(JSON_salts)
    if err != nil{
        return errors.New("unauthorized")
    }

    var salt string

    if keyExists(salts, account){
        salt = salts[account]
    }

    h := hash(password, salt)

    if hashes[account] != h{
        return errors.New("unauthorized")
    }
    return nil
}

func keyExists(decoded map[string]string, key string) bool {
    val, ok := decoded[key]
    return ok && val != ""
}

func unmarshal(jsonFile string) (map[string]string, error) {
    //read json file
    file, err := os.Open(jsonFile)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    defer file.Close()
    byteValue, _ := ioutil.ReadAll(file)

    // unmarshal the data
    var data map[string]string
    err = json.Unmarshal(byteValue, &data)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return data, err
}

func marshal(jsonFile string, data map[string]string) error {
    obj, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
        return err
    }

    fmt.Println(string(obj))

    err = ioutil.WriteFile(jsonFile, obj, 0644)
    return err
}

func keyGen(password string) string {
	//hash the master password to get a sufficient amount of bytes
	bytes := []byte(hash(password, "SOME_SALT"))

	//return string representation of first 32 bytes
	return hex.EncodeToString(bytes[0:32])
}

func getSecret(username string) string {
    secrets, err := unmarshal(JSON_secrets)
    if err != nil{
        panic(err)
    }
    var secret string
    if keyExists(secrets, username){
        secret = secrets[username]
    } else {
        secret = genSecret()
        secrets[username] = secret
        marshal(JSON_secrets, secrets)
    }
    
    return decrypt(secret, keyGen(MASTER_PASSWORD))
}

func genSecret() string {
	//generate random secret //TODO: secure random function?
	secret := make([]byte, 10)
	if _, err := rand.Read(secret); err != nil {
		panic(err)
	}

	//encrypt secret
	encrypted := encrypt(base32.StdEncoding.EncodeToString(secret), keyGen(MASTER_PASSWORD))

	//return the encrypted secret
	return encrypted
}