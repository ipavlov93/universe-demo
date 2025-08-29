package controller

import (
	"fmt"
	"strings"
)

func GetIDFromPath(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid path format")
	}

	return parts[2], nil
}
