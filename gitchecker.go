package main

import (
	"fmt"
	"github.com/gookit/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type summary struct {
	countGits      int
	countUpdated   int
	countUnpushed  int
	countChanges   int
	countUntracked int
}

func Summary(cSum *summary) {

	fmt.Println("Summary: (over the " + strconv.Itoa(cSum.countGits) + " repositories)")
	fmt.Println("\tUpdated " + strconv.Itoa(cSum.countUpdated) + "/" + strconv.Itoa(cSum.countGits))
	fmt.Println("\tUnPushed " + strconv.Itoa(cSum.countUnpushed) + "/" + strconv.Itoa(cSum.countGits))
	fmt.Println("\tChanges not committed " + strconv.Itoa(cSum.countChanges) + "/" + strconv.Itoa(cSum.countGits))
	fmt.Println("\tUntracked files " + strconv.Itoa(cSum.countUntracked) + "/" + strconv.Itoa(cSum.countGits))
}

func AddToChecker(path string, cSum *summary) {
	cmdGit := exec.Command("git", "status")
	cmdGit.Dir = path
	out, err := cmdGit.Output()

	if err != nil {
		log.Fatal(err)
	}

	cSum.countGits++

	st := string(out)
	if strings.Contains(st, "Your branch is up to date with") {
		color.New(color.FgWhite).Println("\t\t* Local changes updated with remotes")
		cSum.countUpdated++
	}
	if strings.Contains(st, "Your branch is ahead of") {
		color.New(color.FgCyan).Println("\t\t* Unpushed commits")
		cSum.countUnpushed++
	}
	if strings.Contains(st, "Changes not staged for commit") {
		color.New(color.FgRed).Println("\t\t* Has changes not staged for commit")
		cSum.countChanges++
	}
	if strings.Contains(st, "Untracked files:") {
		color.New(color.FgYellow).Println("\t\t* Has untracked files")
		cSum.countUntracked++
	}

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
	// root := "/home/sido/Projects/"
	root := "." // Current directory by default

	if len(os.Args) == 2 {
		root = string(os.Args[1])
	}

	fmt.Println("Git full checker!")

	var files []string
	var err error
	var cSum summary

	files, err = GetDirWalk(root)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println("At : " + file)
		AddToChecker(file, &cSum)
	}
	Summary(&cSum)
}
