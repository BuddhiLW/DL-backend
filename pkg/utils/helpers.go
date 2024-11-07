package utils

import "log"

// LogError logs an error message to the console.
func LogError(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}
