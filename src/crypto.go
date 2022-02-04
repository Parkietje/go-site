package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"
)

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
