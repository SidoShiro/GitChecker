package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"shell_color"
	"strconv"
	"strings"
)

type funcPrint func(str string)

type summary struct {
	countGits         int
	countUpdated      int
	countUnpushed     int
	countChanges      int
	countUntracked    int
	listUpdatedGits   []string
	listUnpushedGits  []string
	listChangesGits   []string
	listUntrackedGits []string
}

func printList(printFunc funcPrint, listGits []string) {
	for _, v := range listGits {
		printFunc(" " + v)
	}
}

func Summary(cSum *summary, showDetail bool) {

	shell_color.BluePrintln("Summary: (over the " + strconv.Itoa(cSum.countGits) + " repositories)")
	shell_color.GreenPrint("\tUpdated " + strconv.Itoa(cSum.countUpdated) + "/" + strconv.Itoa(cSum.countGits))
	if showDetail {
		if len(cSum.listUpdatedGits) != 0 {
			fmt.Print("  -:")
			printList(shell_color.GreenPrint, cSum.listUpdatedGits)
		}
	}
	fmt.Print("\n")
	shell_color.CyanPrint("\tUnPushed " + strconv.Itoa(cSum.countUnpushed) + "/" + strconv.Itoa(cSum.countGits))
	if showDetail {
		if len(cSum.listUnpushedGits) != 0 {
			fmt.Print("  -:")
			printList(shell_color.CyanPrint, cSum.listUnpushedGits)
		}
	}
	fmt.Print("\n")

	shell_color.RedPrint("\tChanges not committed " + strconv.Itoa(cSum.countChanges) + "/" + strconv.Itoa(cSum.countGits))
	if showDetail {
		if len(cSum.listChangesGits) != 0 {
			fmt.Print("  -:")
			printList(shell_color.RedPrint, cSum.listChangesGits)
		}
	}
	fmt.Print("\n")

	shell_color.YellowPrint("\tUntracked files " + strconv.Itoa(cSum.countUntracked) + "/" + strconv.Itoa(cSum.countGits))
	if showDetail {
		if len(cSum.listUntrackedGits) != 0 {
			fmt.Print("  -:")
			printList(shell_color.YellowPrint, cSum.listUntrackedGits)
		}
	}
	fmt.Print("\n")
}

func AddToChecker(path string, cSum *summary, optionVerbose bool) {
	cmdGit := exec.Command("git", "status")
	cmdGit.Dir = path
	out, err := cmdGit.Output()

	if err != nil {
		log.Fatal(err)
	}

	cSum.countGits++

	st := string(out)
	if strings.Contains(st, "Your branch is up to date with") {
		cSum.listUpdatedGits = append(cSum.listUpdatedGits, path)
		if optionVerbose {
			shell_color.GreenPrintln("\t\t* Local changes updated with remotes")
		}
		cSum.countUpdated++
	}
	if strings.Contains(st, "Your branch is ahead of") {
		cSum.listUnpushedGits = append(cSum.listUnpushedGits, path)
		if optionVerbose {
			shell_color.CyanPrintln("\t\t* Unpushed commits")
		}
		cSum.countUnpushed++
	}
	if strings.Contains(st, "Changes not staged for commit") {
		cSum.listChangesGits = append(cSum.listChangesGits, path)
		if optionVerbose {
			shell_color.RedPrintln("\t\t* Has changes not staged for commit")
		}
		cSum.countChanges++
	}
	if strings.Contains(st, "Untracked files:") {
		cSum.listUntrackedGits = append(cSum.listUntrackedGits, path)
		if optionVerbose {
			shell_color.YellowPrintln("\t\t* Has untracked files")
		}
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
	verbose := ""
	optionVerbose := false

	if len(os.Args) >= 2 {
		root = string(os.Args[1])
	}
	if len(os.Args) >= 3 {
		verbose = string(os.Args[2])
	}

	if verbose != "" {
		optionVerbose = true
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
		if optionVerbose {
			fmt.Println("At : " + file)
		}
		AddToChecker(file, &cSum, optionVerbose)
	}
	Summary(&cSum, true)
}
