package main

import (
	"flag"
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

type funcPrint func(s string)

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
	shell_color.BluePrint("Summary: (over the " + strconv.Itoa(cSum.countGits) + " repositories)\n")
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
			shell_color.GreenPrint("\t\t* Local changes updated with remotes\n")
		}
		cSum.countUpdated++
	}
	if strings.Contains(st, "Your branch is ahead of") {
		cSum.listUnpushedGits = append(cSum.listUnpushedGits, path)
		if optionVerbose {
			shell_color.CyanPrint("\t\t* Unpushed commits\n")
		}
		cSum.countUnpushed++
	}
	if strings.Contains(st, "Changes not staged for commit") {
		cSum.listChangesGits = append(cSum.listChangesGits, path)
		if optionVerbose {
			shell_color.RedPrint("\t\t* Has changes not staged for commit\n")
		}
		cSum.countChanges++
	}
	if strings.Contains(st, "Untracked files:") {
		cSum.listUntrackedGits = append(cSum.listUntrackedGits, path)
		if optionVerbose {
			shell_color.YellowPrint("\t\t* Has untracked files\n")
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
				// Search for an .git as Dir
				if f.Name() == ".git" && f.IsDir() {
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
	rootPtr := flag.String("root", ".", "Folder, current directory by default (Required)")
	verbosePtr := flag.Bool("verbose", false, "Show detailed information")
	helpPtr := flag.Bool("help", false, "Show this help.")
	flag.Parse()

	if *rootPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *helpPtr == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("textPtr: %s, metricPtr: %t, uniquePtr: %t\n", *rootPtr, *verbosePtr, *helpPtr)

	// root := "/home/sido/Projects/"
	// Current directory by default
	root := *rootPtr
	optionVerbose := *verbosePtr

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
