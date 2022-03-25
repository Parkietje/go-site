package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
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
	} else if len(parts) == 1 { // windows path
		parts = strings.Split(path, "\\")
		filename = parts[len(parts)-1]
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

func listBlobs() (string, error) {
	result := ""

	if STORAGE_ACCOUNT == "" || CONTAINER_NAME == "" || STORAGE_KEY == "" {
		return result, errors.New("No azure credentials supplied")
	}

	cmd := exec.Command("az", "storage", "blob", "list",
		"--account-name", STORAGE_ACCOUNT,
		"--container-name", CONTAINER_NAME,
		"--account-key", STORAGE_KEY,
	)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}
	values := gjson.Get(string(stdout), "#.name")
	println(values.String())

	return values.String(), nil
}

func deployAzureMongo(name string) (string, error) {
	result := ""

	if STORAGE_ACCOUNT == "" || CONTAINER_NAME == "" || STORAGE_KEY == "" {
		return result, errors.New("No azure credentials supplied")
	}

	cmd := exec.Command("az", "vm", "create",
		"-n", name,
		"-g", "RG_UBIOPS",
		"--vnet-name", "VNET-ubiops",
		"--subnet", "mongo-subnet",
		"--image", "ubuntults",
		"--admin-username", "ubuntu",
		"--public-ip-address", "",
		"--nsg", "mongodb-nsg",
		"--os-disk-size-gb", "30",
		"--size", "Standard_B2ms",
		"--ssh-key-values", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC+ghGzHmkkHAjQ6haCim6ssXtAWdrVrzLU8yA2rE4tFEhMxt0R6+31W3KeLBnJR9Mt7uyNlLHBpgURDkPfqLy3WN5HnetoNaA2qBFbEjgT+khu6h0tGllf+PqM4UgrvPYe3HJdUS/VWQzHvnWvG/PvQNrSF+IiduvF4osx+2/+oZ+kOT9Wu0usVUoZRIcgQHtpptul1HTTVMXT8ggj14ywzgnqeYrGwjBOYRqTVKFsJaTSaW8/CCm84tVSZgdS8DSwLVKSXO1uPXdBdjjX2OAhKaGcFsT+yAJhLzWeGgvN1lIcs+SPUuV5MsMYGlAxp3AL/cCprMC9NnSPPkqbdzWp1j8V0a1NFJqXu6oMj4fm/dUESU2yQ9JW0YURB8dncHGpptId5GkOcB/uFP2yrQK2b+2U+Yoi0xlC+AOdu2kBoorHB4DjySJzR8IGEwB/etrq7ZkdiBHA2RQ5nsItSQRJSzU8k4G/m63C2Re1ChBqVMydUZhgpzj803j9ynHIX9k= azuread\\yannichiodi@LAPTOP-NQIP5U8V",
	)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}
	values := gjson.Get(string(stdout), "privateIpAddress")
	println(values.String())

	return values.String(), nil
}

func listVMs() (string, error) {
	//result := []string{}
	result := ""
	cmd := exec.Command("az", "vm", "list",
		"-g", "RG_UBIOPS",
	)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}
	values := gjson.Get(string(stdout), "#.name")
	println(values.String())

	return values.String(), nil
}

func getIP(VM string) (string, error) {
	//result := []string{}
	result := ""
	cmd := exec.Command("az", "vm", "list-ip-addresses",
		"-g", "RG_UBIOPS",
		"-n", VM,
	)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}
	values := gjson.Get(string(stdout), "0.virtualMachine.network.privateIpAddresses.0")

	return values.String(), nil
}

func SCP(folder string, destinationIP string) error {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	files = files[1:]
	for _, fullpath := range files {

		fmt.Println("scp " + "-o StrictHostKeyChecking=no " + fmt.Sprint(fullpath) + " ubuntu@" + fmt.Sprint(destinationIP) + ":/home/ubuntu/")

		cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", fullpath, "ubuntu@"+fmt.Sprint(destinationIP)+":/home/ubuntu/")

		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(string(stdout))

	}

	return nil
}

func execute(command string, destinationIP string) error {

	fmt.Println("ssh " + "-o StrictHostKeyChecking=no " + " ubuntu@" + fmt.Sprint(destinationIP) + " " + command)

	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", "C:/Users/YanniChiodi/.ssh/id_rsa", "ubuntu@"+fmt.Sprint(destinationIP), command)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(stdout))

	return nil
}
