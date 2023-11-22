package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getLatestTag() string {
	// git status
	out, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	} else {
		fmt.Println("Current working directory:", string(out))
	}
	out, err = exec.Command("ls").Output()
	if err != nil {
		fmt.Println("Error listing files:", err)
		os.Exit(1)
	} else {
		fmt.Println("Files:", string(out))
	}
	out, err = exec.Command("which", "git").Output()
	if err != nil {
		fmt.Println("Error finding git:", err)
		os.Exit(1)
	} else {
		fmt.Println("Git location:", string(out))
	}
	out, err = exec.Command("git", "status").Output()
	if err != nil {
		fmt.Println("Error fetching git status:", err)
		os.Exit(1)
	} else {
		fmt.Println("Git status:", string(out))
	}
	out, err = exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
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
		if strings.Contains(file, pattern) {
			fmt.Println("Matched file:", file)
			matchedFiles = append(matchedFiles, file)
		}
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

func runTests(testCmd string) error {
	name := strings.Split(testCmd, " ")[0]
	args := strings.Split(testCmd, " ")[1:]
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var searchPattern string
	var testCmd string

	flag.StringVar(&searchPattern, "pattern", "", "Search pattern for public API tests")
	flag.StringVar(&testCmd, "testCmd", "", "Command to run tests")
	flag.Parse()

	if searchPattern == "" {
		searchPattern = os.Getenv("INPUT_PATTERN")
		fmt.Println("Received pattern:", searchPattern, "from env var")
		if searchPattern == "" {
			fmt.Println("Search pattern must be specified.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Received pattern:", searchPattern, "from flag")
	}

	if testCmd == "" {
		testCmd = os.Getenv("INPUT_TESTCMD")
		fmt.Println("Received testCmd:", testCmd, "from env var")
		if testCmd == "" {
			fmt.Println("Test command must be specified.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Received testCmd:", testCmd, "from flag")
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

	// debug logs
	cwd, _ := os.Getwd()
	fmt.Println("Current Working Directory:", cwd)
	fmt.Println("Environment PATH:", os.Getenv("PATH"))
	cmd := exec.Command("go", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error checking Go version:", err)
		os.Exit(1)
	}
	// end debug logs

	err = runTests(testCmd)
	fmt.Println("Ran tests:", err)
	if err != nil {
		fmt.Println("\nBreaking changes detected!")
		os.Exit(1)
	} else {
		fmt.Println("\nNo breaking changes detected.")
	}
}
