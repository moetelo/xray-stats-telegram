package internal

import (
	"strings"
	"time"
)

func ParseDate(messageText string) (time.Time, error) {
	args := strings.Fields(messageText)
	if len(args) < 2 {
		return time.Now(), nil
	}

	return time.Parse(time.DateOnly, args[1])
}
