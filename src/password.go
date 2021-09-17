package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
	"golang.org/x/term"
)

const (
	JSON_hashes  = "../data/accounts/hashes.json"
	JSON_salts   = "../data/accounts/salts.json"
	JSON_secrets = "../data/accounts/secrets.json"
	STATS_FILE   = "../data/stats/webstats.json"
	IP_FILE      = "../data/stats/ip-addresses.json"
)

var (
	ADMIN = "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
)

func hash(password string, salt string) string {
	h := sha3.New512()
	h.Write([]byte(password + salt))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

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
	}

	h := hash(password, salt)

	if hashes[account] != h {
		fmt.Println("password mismatch")
		return errors.New("unauthorized")
	}
	return nil
}

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

// read file and return decrypted secret for user
func getSecret(username string) string {
	secrets, err := unmarshal(JSON_secrets)
	if err != nil {
		panic(err)
	}
	var secret string
	if keyExists(secrets, username) {
		secret = secrets[username]
	} else {
		secret = genSecret()
		secrets[username] = secret
		marshal(JSON_secrets, secrets)
	}

	return decrypt(secret, keyGen(MASTER_PASSWORD))
}

// generate and encrypt a new random QR secret
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

// derive long key from master pass
func keyGen(password string) string {
	//hash the master password to get a sufficient amount of bytes
	bytes := []byte(hash(password, "SOME_SALT"))
	//return string representation of first 32 bytes
	return hex.EncodeToString(bytes[0:32])
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
// param user should be the hashed username
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

	// delete from files
	err = del(user, JSON_hashes)
	if err != nil {
		return
	}
	err = del(user, JSON_salts)
	if err != nil {
		return
	}
	return del(user, JSON_secrets)
}

// add admin account from stdin inputs
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
// DO NOT LOG/PRINT SECRET
func readPassword() string {
	fmt.Println("Enter password: ")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println(err)
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

//increment web stat
func addStat(stat string) (err error) {
	jmap, err := unmarshal(STATS_FILE)
	if err != nil {
		return
	}
	count := jmap[stat]
	if i, err := strconv.Atoi(count); err == nil {
		i++
		count = strconv.Itoa(i)
	}
	jmap[stat] = count
	return marshal(STATS_FILE, jmap)
}

//increment web stat
func addIP(ip string) (err error) {
	jmap, err := unmarshal(IP_FILE)
	if err != nil {
		return
	}
	if keyExists(jmap, ip) {
		count := jmap[ip]
		if i, err := strconv.Atoi(count); err == nil {
			i++
			count = strconv.Itoa(i)
		}
	} else {
		jmap[ip] = "1"
	}
	return marshal(IP_FILE, jmap)
}
