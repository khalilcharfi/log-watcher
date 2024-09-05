package main

func AssignColor(file string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := colorMap[file]; !exists {
		colorMap[file] = colors[len(colorMap)%len(colors)]
	}
}