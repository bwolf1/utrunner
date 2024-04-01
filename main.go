package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Start the timer
	start := time.Now()

	// Read the config file
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config map[string]interface{}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Read the directories to skip from the config file
	directoriesToSkip := []string{}
	for _, dir := range config["directoriesToSkip"].([]interface{}) {
		directoriesToSkip = append(directoriesToSkip, dir.(string))
	}

	// Set path to the current directory
	path := config["basePath"].(string)
	reportData := []byte{}
	dirCount := -1 // Start at -1 to exclude the root directory

	// Walk through each directory in the path
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, walkErr error) error {
		// Skip directories that do not contain go files
		for _, dir := range directoriesToSkip {
			if path == dir {
				return fs.SkipDir
			}
		}

		// Get searchDepth from config file and convert to integer
		searchDepth, walkErr := strconv.Atoi(config["searchDepth"].(string))
		if walkErr != nil {
			log.Fatal(walkErr)
		}

		if d.IsDir() {
			// Only walk top-level directories
			// We get better output by letting `go test` handle subdirectory traversal
			if strings.Count(path, string(os.PathSeparator)) > searchDepth {
				return fs.SkipDir
			}

			// Run the tests in the current directory and its subdirectories
			cmd := exec.Command("go", "test", "-v", "./...")
			cmd.Dir = path
			output, walkErr := cmd.CombinedOutput()

			// Ignore exits resulting from no go files in the directory
			if walkErr != nil && walkErr.Error() != "exit status 1" {
				fmt.Println(walkErr)
			}

			// Collect and display output for directories with go files only
			if !strings.HasPrefix(string(output), "pattern ./...: directory prefix") {
				fmt.Println(string(output))
				reportData = append(reportData, output...)
			}

			// Count the number of directories walked
			dirCount++
			fmt.Println("Current directory: ", path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Search the reportData for the number of tests passed, failed and skipped
	passed := strings.Count(string(reportData), "PASS")
	failed := strings.Count(string(reportData), "FAIL")
	skipped := strings.Count(string(reportData), "SKIP")

	// Create a formatted report summary
	sep := "-------------------------------------------------"
	summary := fmt.Sprintf(
		"%s\n"+
			"| %-46s |\n%-22s|\n"+
			"| %-22s | %-22d|\n"+
			"| %-22s | %-22d|\n"+
			"| %-22s | %-22d|\n"+
			"| %-22s | %-22d|\n"+
			"| %-22s | %-22d|\n"+
			"| %-22s | %-22d|\n"+
			"%s\n",
		sep,
		"Report",
		sep,
		"Directories walked:",
		dirCount,
		"Tests passed:",
		passed,
		"Tests failed:",
		failed,
		"Tests run:",
		passed+failed,
		"Tests skipped:",
		skipped,
		"Tests written:",
		passed+failed+skipped,
		sep,
	)

	// Print the report summary to the console
	fmt.Println(summary)

	// Write all report data and the report summary to a report file
	reportName := "utrunner-report.txt" // TODO make this configurable
	err = os.WriteFile(reportName, append(reportData, []byte(summary)...), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Print the report file path and the run time to the console
	fmt.Println("Report file path: ", filepath.Join(path, reportName))
	fmt.Println("Run time: ", time.Since(start).Round(time.Second))
}
