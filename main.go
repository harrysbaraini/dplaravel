package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type deployInfo struct {
	Name     string
	Format   string
	Branch   string
	Dest     string
	Commands []command
}

type command struct {
	Cmd string
}

// readJson fetchs the JSON configuration file
func readJSON() []byte {
	file, err := ioutil.ReadFile("./dplaravel.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	return file
}

func dirExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return true, err
	}
	return true, nil
}

// Main function

func main() {
	fmt.Println("######################")
	fmt.Println("Laravel Deploy on SFTP")
	fmt.Println("######################")

	// gets the current working directory and checks if it is a .git directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Working directory: " + dir)

	if _, err := dirExists("./.git"); err != nil {
		log.Fatal(err)
	}

	// reads the dplaravel.json file to retrieve its configuration
	var dpInfo deployInfo
	json.Unmarshal(readJSON(), &dpInfo)
}
