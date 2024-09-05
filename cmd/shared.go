package main

import "sync"

var (
	mu        sync.Mutex
	colorMap  = make(map[string]string)
	colors    = []string{
		"\033[31m", // Red
		"\033[32m", // Green
		"\033[33m", // Yellow
		"\033[34m", // Blue
		"\033[35m", // Magenta
		"\033[36m", // Cyan
	}
	resetColor = "\033[0m"
)