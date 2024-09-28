package logger

import (
	"fmt"
	"time"
)

var debugMode bool

func SetDebugMode(debug bool) {
	debugMode = debug
}

func Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", timestamp, message)
}

func Debug(message string) {
	if debugMode {
		Log(fmt.Sprintf("DEBUG: %s", message))
	}
}

func Info(message string) {
	Log(fmt.Sprintf("INFO: %s", message))
}

func Warning(message string) {
	Log(fmt.Sprintf("WARNING: %s", message))
}

func Error(message string) {
	Log(fmt.Sprintf("ERROR: %s", message))
}
