package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AddToChecker(path string) ([]byte, error) {
	cmdGit := exec.Command("git", "status")
	cmdGit.Dir = path
	out, err := cmdGit.Output()

	if err != nil {
		log.Fatal(err)
	}

	st := string(out)

	if strings.Contains(st, "Your branch is up to date with") {
		fmt.Println("\t\t* Local changes updated with remotes")
	}
	if strings.Contains(st, "Your branch is ahead of") {
		fmt.Println("\t\t* Unpushed commits")
	}
	if strings.Contains(st, "Changes not staged for commit") {
		fmt.Println("\t\t* Has changes not staged for commit")
	}
	if strings.Contains(st, "Untracked files:") {
		fmt.Println("\t\t* Has untracked files")
	}
	// FIXME OLD
	/*
		else {
			fmt.Printf("%s", out)
		}
	*/

	return out, err
}

func GetDirWalk(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name()[0] != '.' && info.Name()[0] != '$' {
			filesInDir, err2 := ioutil.ReadDir(path)
			if err2 != nil {
				log.Fatal(err2)
			}

			for _, f := range filesInDir {
				if f.Name() == ".git" {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})
	return files, err
}

func main() {
	root := "/home/sido/Projects/"
	if len(os.Args) == 2 {
		root = string(os.Args[1])
	}

	fmt.Println("Git full checker!")

	var files []string
	var err error

	files, err = GetDirWalk(root)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println("At : " + file)
		_, _ = AddToChecker(file)
	}
}
