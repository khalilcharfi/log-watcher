package main

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchFile(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create a new watcher for %s: %v", path, err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go readFile(done, watcher, path)  // Assuming readFile is defined elsewhere

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Failed to add %s to watcher: %v", path, err)
	}

	<-done
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

	files, err := FindLogFiles(dir)
	if err != nil {
		log.Fatalf("Failed to find log files in directory %s: %v", dir, err)
	}
	for _, file := range files {
		go WatchFile(file)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if filepath.Ext(event.Name) == ".log" {
					log.Printf("New log file detected: %s\n", event.Name)
					go WatchFile(event.Name)
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