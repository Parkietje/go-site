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

	fmt.Printf("test1")

	fmt.Printf(STORAGE_ACCOUNT)
	if STORAGE_ACCOUNT == "" || CONTAINER_NAME == "" || STORAGE_KEY == "" {
		return errors.New("No azure credentials supplied")
	}

	fmt.Printf("test2")

	var filename string
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		filename = parts[len(parts)-1]
	} else if len(parts) == 1 {
		filename = parts[0]
	} else {
		return errors.New("error parsing filename from path: \n" + path)
	}
	fmt.Printf("test3")

	cmd := exec.Command("az", "storage", "blob", "upload",
		"--account-name", STORAGE_ACCOUNT,
		"--container-name", CONTAINER_NAME,
		"--account-key", STORAGE_KEY,
		"--name", filename,
		"--file", path,
	)
	stdout, err := cmd.Output()

	fmt.Printf("test4")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Print the output
	fmt.Println(string(stdout))
	return nil
}
