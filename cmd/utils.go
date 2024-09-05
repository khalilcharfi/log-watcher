package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func IsFile(file string) bool {
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
	// This function parses a filter string in the format key[value]
	re := regexp.MustCompile(`^(\w+)$begin:math:display$(.+)$end:math:display$$`)
	matches := re.FindStringSubmatch(filter)
	if len(matches) != 3 {
		return fmt.Errorf("filter must be in the format key[value]")
	}
	filterKey = matches[1]
	filterValue = matches[2]
	return nil
}

func readFile(done chan bool, watcher *fsnotify.Watcher, path string) {
	defer close(done)

	// Assign a color to this file if prettify is enabled
	AssignColor(path)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				file, err := os.ReadFile(path)
				if err != nil {
					log.Printf("Failed to read %s: %v", path, err)
					continue
				}
				lines := strings.Split(string(file), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						PrintWithColor(path, line)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching %s: %v", path, err)
		}
	}
}