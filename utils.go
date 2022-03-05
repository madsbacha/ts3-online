package main

import (
	"os"
	"strings"
)

func excludeUsername(username string) bool {
	exclude_usernames, exists := os.LookupEnv("EXCLUDE_USERNAMES")
	if !exists {
		return false
	}
	usernames := strings.Split(exclude_usernames, ",")

	for _, val := range usernames {
		if username == val {
			return true
		}
	}

	return false
}
