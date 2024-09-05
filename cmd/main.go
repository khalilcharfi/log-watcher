package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	prettify   bool
	filterKey  string
	filterValue string
)

func main() {
	// Define flags
	dirPath := flag.String("d", "", "Specify the directory to watch for log files")
	flag.BoolVar(&prettify, "p", false, "Enable color-coded output for each file")
	filter := flag.String("filter", "", "Filter logs by a specific key-value pair, e.g., category[debug] or log_level[debug]")
	flag.Parse()

	// Handle command-line arguments
	args := flag.Args()
	if len(args) == 1 {
		file := args[0]
		if !IsFile(file) {
			log.Fatalf("The provided path is not a valid file: %s", file)
		}
		WatchFile(file)
		return
	}

	// Handle directory watching
	if *dirPath == "" && len(args) == 0 {
		log.Fatal("Please provide either a directory to watch using the -d flag or a file path directly.")
	}

	if *filter != "" {
		if err := ParseFilter(*filter); err != nil {
			log.Fatalf("Invalid filter format: %v", err)
		}
	}

	// Start watching the directory if provided
	if *dirPath != "" {
		go WatchDirectory(*dirPath)
		select {} // Keep the program running indefinitely
	}
}