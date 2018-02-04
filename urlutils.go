package main

import (
	"errors"
	"strings"
)

const armEndpoint string = "https://management.azure.com"

func isArmURLPath(urlPath string) bool {
	return strings.HasPrefix(urlPath, "/subscriptions") ||
		strings.HasPrefix(urlPath, "/tenants") ||
		strings.HasPrefix(urlPath, "/providers")
}

func getRequestURL(path string) (string, error) {
	if !isArmURLPath(path) {
		return "", errors.New("Url path specified is invalid")
	}

	return armEndpoint + path, nil
}
