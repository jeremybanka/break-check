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

	fmt.Println("Filtering files by pattern:", pattern)
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
		fmt.Println("Matched file:", string(out))
		matchedFiles = append(matchedFiles, strings.TrimSpace(string(out)))
	}
	return matchedFiles, nil
}

func checkoutFromTag(tag string, files []string) error {
	fmt.Println("Checking out files at tag:", tag)
	cmdArgs := []string{"checkout", tag, "--"}
	cmdArgs = append(cmdArgs, files...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runTests(testCmd string, files []string) error {
	fmt.Println("Running tests:", testCmd, files)
	cmdArgs := strings.Fields(testCmd)
	cmdArgs = append(cmdArgs, "--")
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
		fmt.Println("No tests matched this pattern:", searchPattern)
		os.Exit(2)
	}

	err = checkoutFromTag(tag, matchedFiles)
	if err != nil {
		fmt.Println("Error checking out files at tag:", err)
		os.Exit(1)
	}

	err = runTests(testCmd, matchedFiles)
	fmt.Println("Ran tests:", err)
	if err != nil {
		fmt.Println("\nBreaking changes detected!")
		os.Exit(1)
	} else {
		fmt.Println("\nNo breaking changes detected.")
	}
}
