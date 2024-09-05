package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func IsFile(file string) bool {
	return fileExists(file)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FindLogFiles(dir string) ([]string, error) {
	var logFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".log" {
			logFiles = append(logFiles, path)
		}
		return nil
	})
	return logFiles, err
}

func ParseFilter(filter string) error {
	re := regexp.MustCompile(`^(\w+)$begin:math:display$(.+)$end:math:display$$`)
	matches := re.FindStringSubmatch(filter)
	if len(matches) != 3 {
		return fmt.Errorf("filter must be in the format key[value]")
	}
	filterKey = matches[1]
	filterValue = matches[2]
	return nil
}