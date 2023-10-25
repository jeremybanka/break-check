package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var testDir string

// setupTestEnvironment sets up the test environment
func setupTestEnvironment(t *testing.T) {
	// Create a temporary directory
	var err error
	testDir, err = os.MkdirTemp("", "break-check-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}

	buildCmd := exec.Command("go", "build", "-o", "break-check")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build break-check: %s", err)
	}

	// Copy the binary to the test environment
	os.Rename("break-check", filepath.Join(testDir, "break-check"))

	// Change directory to the test environment
	os.Chdir(testDir)

	// Initialize a Git repository
	exec.Command("git", "init").Run()

	// Set up Node.js environment
	exec.Command("npm", "init", "-y").Run()

	// Write the source file
	srcContent := `
class MyClass {
  constructor() {
    this.privateVar = "private";
  }

  publicMethod() {
    return "publicMethodOutput";
  }

  privateMethod() {
    return this.privateVar;
  }
}

module.exports = MyClass;
`
	os.WriteFile(filepath.Join(testDir, "src.js"), []byte(srcContent), 0644)

	// Write the public API test file
	publicTestContent := `
const assert = require('assert').strict;
const MyClass = require('./src');

describe('MyClass Public API Test', function() {
  it('should test publicMethod', function() {
    const obj = new MyClass();
    assert.strictEqual(obj.publicMethod(), "publicMethodOutput");
  });
});
`
	os.WriteFile(filepath.Join(testDir, "publicTest.js"), []byte(publicTestContent), 0644)

	// Write the private API test file
	privateTestContent := `
const assert = require('assert').strict;
const MyClass = require('./src');

describe('MyClass Private API Test', function() {
  it('should test privateMethod', function() {
    const obj = new MyClass();
    assert.strictEqual(obj.privateMethod(), "private");
  });
});
`
	os.WriteFile(filepath.Join(testDir, "privateTest.js"), []byte(privateTestContent), 0644)

	// Commit and tag the initial state
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial commit with public and private tests").Run()
	exec.Command("git", "tag", "v1.0").Run()
}

// tearDownTestEnvironment cleans up the test environment
func tearDownTestEnvironment(t *testing.T) {
	os.RemoveAll(testDir) // remove the temporary directory
}

func TestBreakCheck(t *testing.T) {
	setupTestEnvironment(t)
	defer tearDownTestEnvironment(t)

	// Introduce a breaking change to src.js
	srcFilePath := filepath.Join(testDir, "src.js")
	srcContent, err := os.ReadFile(srcFilePath)
	if err != nil {
		t.Fatalf("Failed to read src.js: %s", err)
	}

	modifiedSrcContent := strings.Replace(string(srcContent), `"publicMethodOutput"`, `"modifiedPublicMethodOutput"`, 1)
	err = os.WriteFile(srcFilePath, []byte(modifiedSrcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write modified src.js: %s", err)
	}

	// Run the break-check tool
	cmd := exec.Command("./break-check", "--pattern", "Public API Test", "--testCmd", "npm test")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("Expected break-check to report a breaking change, but it didn't.\nOutput: %s", output)
	}

	if !strings.Contains(string(output), "Breaking changes detected!") {
		t.Errorf("Expected 'Breaking changes detected!' in output but got: %s", output)
	}
}
