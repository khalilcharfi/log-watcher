package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	prettify   bool
	colorMap   map[string]string
	mu         sync.Mutex
	colors     = []string{
		"\033[31m", // Red
		"\033[32m", // Green
		"\033[33m", // Yellow
		"\033[34m", // Blue
		"\033[35m", // Magenta
		"\033[36m", // Cyan
	}
	resetColor   = "\033[0m"
	filterKey    string
	filterValue  string
)

func main() {
	// Define the -d flag for directory input, -p flag for prettify, and --filter flag
	dirPath := flag.String("d", "", "Specify the directory to watch for log files")
	flag.BoolVar(&prettify, "p", false, "Enable color-coded output for each file")
	filter := flag.String("filter", "", "Filter logs by a specific key-value pair, e.g., category[debug] or log_level[debug]")
	flag.Parse()

	// Handle the case where a single file path is passed directly as a command-line argument
	args := flag.Args()
	if len(args) == 1 {
		file := args[0]
		if !IsFile(file) {
			log.Fatalf("The provided path is not a valid file: %s", file)
		}
		Watch(file)
		return
	}

	// Handle the case where directory watching and filtering is used
	if *dirPath == "" && len(args) == 0 {
		log.Fatal("Please provide either a directory to watch using the -d flag or a file path directly.")
	}

	if *filter != "" {
		if err := parseFilter(*filter); err != nil {
			log.Fatalf("Invalid filter format: %v", err)
		}
	}

	colorMap = make(map[string]string)

	// Start watching the directory for .log files if a directory is provided
	if *dirPath != "" {
		go WatchDirectory(*dirPath)
		select {} // Keep the program running indefinitely
	}
}

func parseFilter(filter string) error {
	re := regexp.MustCompile(`^(\w+)$begin:math:display$(.+)$end:math:display$$`)
	matches := re.FindStringSubmatch(filter)
	if len(matches) != 3 {
		return fmt.Errorf("filter must be in the format key[value]")
	}
	filterKey = matches[1]
	filterValue = matches[2]
	return nil
}

func IsFile(file string) bool {
	return fileExists(file)
}

func FileInfo(file string) {
	info, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return
	}
	mode := info.Mode()
	fmt.Printf("\n File: %s \t \t Mode: %s \n \n ", file, mode)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Watch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create a new watcher for %s: %v", path, err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go readFile(done, watcher, path)

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Failed to add %s to watcher: %v", path, err)
	}

	<-done
}

func readFile(done chan bool, watcher *fsnotify.Watcher, path string) {
	defer close(done)
	assignColor(path)

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
						printWithColor(path, line)
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

func WatchDirectory(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create directory watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatalf("Failed to add directory %s to watcher: %v", dir, err)
	}

	// Watch the existing .log files in the directory
	files, err := findLogFiles(dir)
	if err != nil {
		log.Fatalf("Failed to find log files in directory %s: %v", dir, err)
	}
	for _, file := range files {
		go Watch(file)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if filepath.Ext(event.Name) == ".log" {
					fmt.Printf("New log file detected: %s\n", event.Name)
					go Watch(event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching directory %s: %v", dir, err)
		}
	}
}

func findLogFiles(dir string) ([]string, error) {
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

func assignColor(file string) {
	if prettify {
		mu.Lock()
		defer mu.Unlock()
		if _, exists := colorMap[file]; !exists {
			colorMap[file] = colors[len(colorMap)%len(colors)]
		}
	}
}

func printWithColor(file, text string) {
	if prettify {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(text), &logEntry); err == nil {
			// Apply the filter if specified
			if filterKey != "" && logEntry[filterKey] != filterValue {
				return
			}

			formattedLog := formatLogEntry(logEntry)
			mu.Lock()
			color := colorMap[file]
			mu.Unlock()
			fmt.Printf("%s%s%s\n", color, formattedLog, resetColor)
		} else {
			// If the line is not a valid JSON, print it as is
			mu.Lock()
			color := colorMap[file]
			mu.Unlock()
			fmt.Printf("%s%s%s\n", color, text, resetColor)
		}
	} else {
		fmt.Println(text)
	}
}

func formatLogEntry(logEntry map[string]interface{}) string {
	// Format the JSON log entry to a more readable format
	timestamp := logEntry["@timestamp"]
	level := logEntry["log_level"]
	app := logEntry["app_name"]
	message := logEntry["message"]
	environment := logEntry["environment"]

	return fmt.Sprintf("[%s] [%s] [%s] %s (%s)",
		timestamp, level, app, message, environment)
}