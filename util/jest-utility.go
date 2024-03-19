package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileCoverage struct {
	FileShortPath string `json:"fileShortPath"`
	CoveredLines  int    `json:"coveredLines"`
}

type CoverageSummary struct {
	Data []FileCoverage `json:"data"`
}

func findRelatedSourceFiles(testFile string, rootDir string) ([]string, error) {
	cmd := exec.Command("npx", "jest", testFile, "--json", "--coverage","--coverageFile=coverage/coverage.json" "--coverageReporters=json-summary")
	cmd.Dir = rootDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Jest error output:", stderr.String())
		return nil, fmt.Errorf("failed to run Jest: %v", err)
	}
	for x, _ := range output {
		fmt.Println(string(rune(x)))
	}

	var coverageSummary CoverageSummary
	err = json.Unmarshal(output, &coverageSummary)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage summary: %v", err)
	}

	affectedSourceFiles := make([]string, 0)
	for _, fileCoverage := range coverageSummary.Data {
		if fileCoverage.CoveredLines > 0 {
			affectedSourceFiles = append(affectedSourceFiles, filepath.Join(rootDir, fileCoverage.FileShortPath))
		}
	}

	return affectedSourceFiles, nil
}

func main() {
	testFileOrPattern := "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing/src/foo.test.js"
	rootDir := "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing"

	affectedSourceFiles, err := findRelatedSourceFiles(testFileOrPattern, rootDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Source files affected by", testFileOrPattern)
	fmt.Println(strings.Join(affectedSourceFiles, "\n"))
}
