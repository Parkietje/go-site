package main

import (
	"bufio"
	"encoding/base64"
	"net/url"
	"os"

	"github.com/dgryski/dgoogauth"
	"rsc.io/qr"
)

func genQR(account string, secret string) string {
	issuer := "webportal"
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}
	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)
	params := url.Values{}
	params.Add("secret", secret)
	params.Add("issuer", issuer)
	URL.RawQuery = params.Encode()
	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(code.PNG())
}

func verify(token string, secret string) (dgoogauth.OTPConfig, error, bool) {
	// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
	// so allocate it once and reuse it for multiple calls.
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// REMOVE THIS BLOCK AFTER TESTING
	if token == "" {
		return *otpc, nil, true
	}
	// REMOVE THIS BLOCK AFTER TESTING

	val, err := otpc.Authenticate(token)
	return *otpc, err, val
}

//encode PNG to html-embeddable string
func imgBase64Str(fileName string) (string, error) {
	imgFile, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	return base64.StdEncoding.EncodeToString(buf), nil
}
