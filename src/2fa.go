package main

import (
	"encoding/base32"
	"io/ioutil"
	"net/url"
	"os"
	"bufio"
	dgoogauth "github.com/dgryski/dgoogauth"
	qr "rsc.io/qr"
	"encoding/base64"
)

const (
	qrFilename = "./ui/static/img/qr.png"
	penguinFilename = "./ui/static/img/pngegg.png"
)

var (
	// Example secret from here:
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	secret = []byte{'H', 'e', 'l', 'l', 'o', '!', 0xDE, 0xAD, 0xBE, 0xEF}
	secretBase32 = base32.StdEncoding.EncodeToString(secret)
)

func genQR(account string) string {
	issuer := "webportal"
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)

	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()
	//fmt.Printf("URL is %s\n", URL.String())

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		panic(err)
	}
	b := code.PNG()
	err = ioutil.WriteFile(qrFilename, b, 0600)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("QR code is in %s. Please scan it into Google Authenticator app.\n", qrFilename)

  	// Embed into an html without PNG file
  	s, err := imgBase64Str(qrFilename)
	if err != nil{
		panic(err)
	}
	return s
}
	
func verify(token string) (dgoogauth.OTPConfig, error, bool) {
	// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
	// so allocate it once and reuse it for multiple calls.
	otpc := &dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(token)

	return *otpc, err, val
}

func imgBase64Str(fileName string) (string , error) {
	imgFile, err := os.Open(fileName) // a QR code image

  	defer imgFile.Close()

  	// create a new buffer base on file size
  	fInfo, _ := imgFile.Stat()
  	var size int64 = fInfo.Size()
  	buf := make([]byte, size)

  	// read file content into buffer
  	fReader := bufio.NewReader(imgFile)
  	fReader.Read(buf)
	
	// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

  	// png.Encode(&buf, image)

  	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
  	return base64.StdEncoding.EncodeToString(buf), err
}