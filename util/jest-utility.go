package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

type TestResults struct {
	TestResults []TestResult `json:"testResults"`
}
type TestResult struct {
	Name string `json:"name"`
}

type CoverageMetrics struct {
	Total   int     `json:"total"`
	Covered int     `json:"covered"`
	Skipped int     `json:"skipped"`
	Percent float64 `json:"pct"`
}

type CoverageData struct {
	Total map[string]CoverageMetrics `json:"total"`
	Files map[string]CoverageMetrics `json:""`
}

type FileCoverage struct {
	FileShortPath string `json:"fileShortPath"`
	CoveredLines  int    `json:"coveredLines"`
}

type CoverageSummary struct {
	Data []FileCoverage `json:"data"`
}
type TestFile struct {
	RootDir      string
	Path         string
	Dependencies []string
	Content      string
	coverageData CoverageData
}

func (t *TestFile) GetCoverageData() ([]string, error) {
	t.GetRelatedComponents()
	cmd := exec.Command("jest",t.Path, "--json", "--coverage", "--outputFile=coverage/coverage.json", "--coverageReporters=json-summary", "--collectCoverageFrom='[" +
		t.GetStringDependencies()+"]'")
	cmd.Dir = t.RootDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Jest error output:", stderr.String())
		return nil, fmt.Errorf("failed to run Jest: %v", err)
	}

	jsonData, err := ioutil.ReadFile(rootDir + "/coverage/coverage.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var testResults TestResults
	if err := json.Unmarshal(jsonData, &testResults); err != nil {
		log.Fatalf("Error unmarshaling JSON data: %v", err)
	}
	var affectedFiles []string
	for _, result := range testResults.TestResults {
		fmt.Println("File Path:", result.Name)
		affectedFiles = append(affectedFiles, result.Name)
	}
	return affectedFiles, nil
}

func (t *TestFile) GetStringDependencies() string {
	projectDir= t.RootDir+""
	var output string
	for _, dependency:= range t.Dependencies{
		output+=
	}

}

//func main() {
//	testFileOrPattern := "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing/src/foo.test.js"
//	rootDir := "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing"
//
//	affectedSourceFiles, err := findRelatedSourceFiles(testFileOrPattern, rootDir)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	fmt.Println("Source files affected by", testFileOrPattern)
//	fmt.Println(strings.Join(affectedSourceFiles, "\n"))
//}
