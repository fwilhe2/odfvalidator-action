package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type LogEntry struct {
	SubFile  string `json:"sub_file"` // file inside the ODS archive
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	CodeLine string `json:"code_line,omitempty"`
	CaretPos int    `json:"caret_pos,omitempty"`
}

type FileReport struct {
	ODSPath string     `json:"ods_path"`
	Entries []LogEntry `json:"entries"`
}

var headerPattern = regexp.MustCompile(`^(?P<path>.+?)(?:\[(?P<line>\d+),(?P<col>\d+)\])?:\s+(?P<severity>Error|Info):\s+(?P<msg>.+)$`)

func main() {
	file, err := os.Open("/odf-errors.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Map of ODS file path → list of entries
	grouped := make(map[string][]LogEntry)

	var currentODS string
	var currentEntry *LogEntry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if m := headerPattern.FindStringSubmatch(line); m != nil {
			// Save previous entry if it exists
			if currentEntry != nil && currentODS != "" {
				grouped[currentODS] = append(grouped[currentODS], *currentEntry)
			}

			fullPath := strings.TrimSpace(m[1])
			odsPath, subFile := splitODSPath(fullPath)

			currentODS = odsPath
			currentEntry = &LogEntry{
				SubFile:  subFile,
				Severity: m[4],
				Message:  strings.TrimSpace(m[5]),
			}
			if m[2] != "" {
				fmt.Sscanf(m[2], "%d", &currentEntry.Line)
				fmt.Sscanf(m[3], "%d", &currentEntry.Column)
			}

		} else if currentEntry != nil && strings.Contains(line, "----^") {
			// Caret position
			currentEntry.CaretPos = strings.Index(line, "^")
		} else if currentEntry != nil && strings.TrimSpace(line) != "" {
			// Code snippet
			currentEntry.CodeLine = strings.TrimSpace(line)
		}
	}

	// Save the last entry
	if currentEntry != nil && currentODS != "" {
		grouped[currentODS] = append(grouped[currentODS], *currentEntry)
	}

	// Convert to JSON
	reports := []FileReport{}
	for odsPath, entries := range grouped {
		reports = append(reports, FileReport{
			ODSPath: odsPath,
			Entries: entries,
		})
	}

	for _, report := range reports {
		for _, e := range report.Entries {
			// GitHub expects repo-relative paths
			repoPath := filepath.Join(report.ODSPath, e.SubFile)
			repoPath = strings.Replace(repoPath, "/github/workspace/", "", 0)
			// Use "warning" or "error" based on severity
			fmt.Printf("::%s file=%s,line=%d,col=%d::%s\n",
				strings.ToLower(e.Severity),
				repoPath,
				e.Line,
				e.Column,
				e.Message,
			)
		}
	}
}

// splitODSPath separates the ODS archive path from the internal file path
func splitODSPath(full string) (odsPath, subFile string) {
	// Example:
	// "/usr/src/./data/common-data-types-de_DE.UTF-8.ods/content.xml"
	// → odsPath: "/usr/src/./data/common-data-types-de_DE.UTF-8.ods"
	//   subFile: "content.xml"
	parts := strings.SplitN(full, ".ods/", 2)
	if len(parts) == 2 {
		return parts[0] + ".ods", parts[1]
	}
	// Handle case where no subpath (just the ODS itself)
	if strings.HasSuffix(full, ".ods") {
		return full, ""
	}
	// If path is weird, just return as-is
	return full, ""
}
