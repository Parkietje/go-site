package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	STORAGE_ACCOUNT string
	CONTAINER_NAME  string
	STORAGE_KEY     string
)

func uploadBlob(path string) error {

	if STORAGE_ACCOUNT == "" || CONTAINER_NAME == "" || STORAGE_KEY == "" {
		return errors.New("No azure credentials supplied")
	}

	var filename string
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		filename = parts[len(parts)-1]
	} else if len(parts) == 1 {
		filename = parts[0]
	} else {
		return errors.New("error parsing filename from path: \n" + path)
	}

	cmd := exec.Command("az", "storage", "blob", "upload",
		"--account-name", STORAGE_ACCOUNT,
		"--container-name", CONTAINER_NAME,
		"--account-key", STORAGE_KEY,
		"--name", filename,
		"--file", path,
	)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Print the output
	fmt.Printf("File upload successful")

	fmt.Println(string(stdout))
	return nil
}
