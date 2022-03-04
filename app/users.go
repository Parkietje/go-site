package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/sha3"
)

// check if provided password matches, return error if not
func passwordCheck(account string, password string) error {
	hashes, err := unmarshal(HASHES)
	if err != nil {
		return errors.New("unauthorized")
	}

	salts, err := unmarshal(SALTS)
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
	err = add(user, salt, SALTS)
	if err != nil {
		return
	}
	err = add(user, hash(password, salt), HASHES)
	if err != nil {
		return
	}
	return add(user, genSecret(), SECRETS)
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
	err = del(user, HASHES)
	if err != nil {
		return err
	}
	err = del(user, SALTS)
	if err != nil {
		return err
	}
	return del(user, SECRETS)
}

// check if key exists in map
func keyExists(decoded map[string]string, key string) bool {
	val, ok := decoded[key]
	return ok && val != ""
}

// read file and return decrypted secret for user
func getSecret(username string) string {
	secrets, err := unmarshal(SECRETS)
	if err != nil {
		panic(err)
	}
	var secret string
	if keyExists(secrets, username) {
		secret = secrets[username]
	} else {
		secret = genSecret()
		secrets[username] = secret
		marshal(SECRETS, secrets)
	}

	return decrypt(secret, keyGen(MASTER_PASSWORD))
}

// generate and encrypt a new random secret
func genSecret() string {
	//generate random secret
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

// hash the given input + salt (provide empty string for no salt)
func hash(text string, salt string) string {
	h := sha3.New512()
	h.Write([]byte(text + salt))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

//borrowed from:
//https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes
func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err.Error())
		return
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

//borrowed from:
//https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes
func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return fmt.Sprintf("%s", plaintext)
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
