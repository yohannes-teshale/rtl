package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func (t *TestFile) GetRelatedComponents() {
	//rootDir := "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing"
	//testFilePath := rootDir + "/src/foo.test.ts"
	projectRoot := t.RootDir + "/src"
	nodeScriptPath := t.RootDir + "/jest-resolve-dependencies.js"

	cmd := exec.Command("node", nodeScriptPath, t.Path, projectRoot)
	cmd.Dir = projectRoot

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %s\n", err)
		fmt.Fprintf(os.Stderr, "Node.js script stderr:\n%s\n", stderr.String())
		return
	}
	var dependencies []string
	err = json.Unmarshal(out.Bytes(), &dependencies)
	if err != nil {
		fmt.Printf("Failed to parse JSON output: %s\n", err)
		return
	}

	t.Dependencies = dependencies
}
