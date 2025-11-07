package utils

import (
	"log"
	"strconv"
)

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid COLLECTOR_PORT: %s", s)
	}
	return i
}
