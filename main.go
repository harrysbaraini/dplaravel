package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type deployInfo struct {
	Name     string
	Branch   string
	Dest     string
	Commands []command
	Sftp     sftpConfig
}

type sftpConfig struct {
	Host      string
	User      string
	Password  string
	Directory string
	Port      string
}

type command struct {
	Cmd string
}

func (c *command) run() {
	a := strings.Split(c.Cmd, " ")
	fmt.Printf(" CMD :: Running %s\n\n", c.Cmd)
	out, err := exec.Command(a[0], a[1:]...).Output()
	if err != nil {
		fmt.Println("Cannot run commmand " + c.Cmd)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
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

// dirExists checks if a directory exists
func dirExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return true, err
	}
	return true, nil
}

// archive creates an archive file from .git repository
func archiveGit(dp deployInfo) string {
	args := []string{"archive"}

	if dp.Dest != "" {
		args = append(args, "--prefix="+dp.Dest+"/")
	} else {
		args = append(args, "--prefix=build/")
	}

	ts := time.Now().UTC()
	output := dp.Name + "." + strings.Replace(dp.Branch, "/", "-", -1) + ts.Format("20060102150405") + ".tar.gz"
	args = append(args, []string{"--format=tar.gz", "--output=" + output, dp.Branch}...)
	fmt.Printf("%v\n", args)

	fmt.Println("Archiving files...")
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
	fmt.Printf("Archive was successfully created[ %s ]\n", output)
	return output
}

// unpack extracts the generated archive to Dest directory
func unpack(file string) {
	out, err := exec.Command("tar", "-xvzf", file).Output()
	if err != nil {
		fmt.Println("Cannot unpack file " + file)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

// sendFiles sends processed directory to SFTP
func sendFiles(dpInfo deployInfo) {
	host := dpInfo.Sftp.User + "@" + dpInfo.Sftp.Host
	args := []string{host, "<<<", "$'put -r " + dpInfo.Dest + "/* " + dpInfo.Sftp.Directory + "'"}

	cmd := exec.Command("sftp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Cannot send files to SFTP")
		log.Fatal(err)
	}
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

	// creates a archive file from .git repository and unpacks it on destination
	f := archiveGit(dpInfo)
	unpack(f)

	// Run any issued command
	os.Chdir(dpInfo.Dest)
	for i := 0; i < len(dpInfo.Commands); i++ {
		dpInfo.Commands[i].run()
	}

	// Clean up distribution
	os.RemoveAll("node_modules/")

	// Send files to SFTP
	// sendFiles(dpInfo)
	final := dpInfo.Name + "." + strings.Replace(dpInfo.Branch, "/", "-", -1) + "-BUILD.tar.gz"
	tarerr := exec.Command("tar", []string{"-czf", final, "."}...).Run()
	if tarerr != nil {
		log.Fatal(tarerr)
	}

	// Clean up generated files
	os.Chdir("..")
	os.Remove(f)
	// os.RemoveAll(dpInfo.Dest)
	fmt.Println("THAT'S OK")
}
