package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	JSON_hashes  = "../data/accounts/hashes.json"
	JSON_salts   = "../data/accounts/salts.json"
	JSON_secrets = "../data/accounts/secrets.json"
)

// check if provided password matches, return error if not
func passwordCheck(account string, password string) error {
	hashes, err := unmarshal(JSON_hashes)
	if err != nil {
		return errors.New("unauthorized")
	}

	salts, err := unmarshal(JSON_salts)
	if err != nil {
		return errors.New("unauthorized")
	}

	var salt string

	if keyExists(salts, account) {
		salt = salts[account]
	} else {
		fmt.Println("user not found")
		return errors.New("unauthorized")
	}

	h := hash(password, salt)

	if hashes[account] != h {
		fmt.Println("password mismatch")
		return errors.New("unauthorized")
	}
	return nil
}

// store new user in /data/accounts files
func addUser(user string, password string, salt string) (err error) {
	user = hash(user, "")

	// add to file lambda
	add := func(key string, value string, file string) (err error) {
		jmap, err := unmarshal(file)
		if err != nil {
			return
		}
		jmap[key] = value
		return marshal(file, jmap)
	}

	//append to files
	err = add(user, salt, JSON_salts)
	if err != nil {
		return
	}
	err = add(user, hash(password, salt), JSON_hashes)
	if err != nil {
		return
	}
	return add(user, genSecret(), JSON_secrets)
}

// delete user info from /data/accounts files
func deleteUser(user string) (err error) {
	// delete lambda
	del := func(user string, file string) error {
		jmap, err := unmarshal(file)
		if err != nil {
			return err
		}
		_, ok := jmap[user]
		if ok {
			delete(jmap, user)
		}
		return marshal(file, jmap)
	}

	user = hash(user, "")

	// delete from files
	err = del(user, JSON_hashes)
	if err != nil {
		return err
	}
	err = del(user, JSON_salts)
	if err != nil {
		return err
	}
	return del(user, JSON_secrets)
}

// add admin account from stdin
func addAdmin() error {
	fmt.Println("Enter admin account name: ")
	var user string
	fmt.Scanln(&user)

	pw := readPassword()

	fmt.Println("Enter salt for admin: ")
	var salt string
	fmt.Scanln(&salt)

	ADMIN = user
	return addUser(user, pw, salt)
}

// read credentials from stdin without echo'ing them in terminal history
func readPassword() string {
	fmt.Println("Enter password: ")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println(err)
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

// check if key exists in map
func keyExists(decoded map[string]string, key string) bool {
	val, ok := decoded[key]
	return ok && val != ""
}

// read json file into go map
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

// write go map to file
func marshal(jsonFile string, data map[string]string) error {
	obj, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return ioutil.WriteFile(jsonFile, obj, 0644)
}
