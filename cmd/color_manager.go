package main

import "sync"

func AssignColor(file string) {
	if prettify {
		mu.Lock()
		defer mu.Unlock()
		if _, exists := colorMap[file]; !exists {
			colorMap[file] = colors[len(colorMap)%len(colors)]
		}
	}
}