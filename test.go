package main

import (
	"os"
	"os/exec"
	"path/filepath"
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

	// TODO: Test the break-check tool
}
