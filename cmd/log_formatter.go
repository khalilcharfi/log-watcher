package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

func PrintWithColor(file, text string) {
	if prettify {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(text), &logEntry); err == nil {
			// Apply the filter if specified
			if filterKey != "" && logEntry[filterKey] != filterValue {
				return
			}

			formattedLog := FormatLogEntry(logEntry)
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

func FormatLogEntry(logEntry map[string]interface{}) string {
	timestamp := logEntry["@timestamp"]
	level := logEntry["log_level"]
	app := logEntry["app_name"]
	message := logEntry["message"]
	environment := logEntry["environment"]

	return fmt.Sprintf("[%s] [%s] [%s] %s (%s)",
		timestamp, level, app, message, environment)
}