package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getLatestTag() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		fmt.Println("Error fetching latest tag:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}

func getTestsAtTag(tag string) ([]string, error) {
	out, err := exec.Command("git", "ls-tree", "-r", tag, "--full-name", "--name-only").Output()
	if err != nil {
		return nil, err
	}
	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	return files, nil
}

func filterFilesByPattern(files []string, pattern string) ([]string, error) {
	var matchedFiles []string

	for _, file := range files {
		fmt.Println("Checking file:", file)
		out, err := exec.Command("grep", "-l", pattern, file).Output()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
				// grep didn't find a match, but that's okay.
				continue
			}
			// An actual error occurred
			fmt.Println("Error running grep:", err)
			return nil, err
		}
		matchedFiles = append(matchedFiles, strings.TrimSpace(string(out)))
	}
	return matchedFiles, nil
}

func runTests(testCmd string, files []string) error {
	cmdArgs := strings.Fields(testCmd)
	cmdArgs = append(cmdArgs, files...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var searchPattern string
	var testCmd string

	flag.StringVar(&searchPattern, "pattern", "", "Search pattern for public API tests")
	flag.StringVar(&testCmd, "testCmd", "go test", "Command to run tests")
	flag.Parse()

	if searchPattern == "" {
		fmt.Println("Search pattern must be specified.")
		os.Exit(1)
	}

	tag := getLatestTag()
	testFiles, err := getTestsAtTag(tag)
	if err != nil {
		fmt.Println("Error fetching test files at tag:", err)
		os.Exit(1)
	}

	matchedFiles, err := filterFilesByPattern(testFiles, searchPattern)
	if err != nil {
		fmt.Println("Error filtering files by pattern:", err)
		os.Exit(1)
	}

	if len(matchedFiles) == 0 {
		fmt.Println("No tests match the specified pattern.")
		os.Exit(1)
	}

	err = runTests(testCmd, matchedFiles)
	if err != nil {
		fmt.Println("\nBreaking changes detected!")
		os.Exit(1)
	} else {
		fmt.Println("\nNo breaking changes detected.")
	}
}
