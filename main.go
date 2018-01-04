package main

import (
	"github.com/tomkukral/cloud-cleaner/cloudformation"
	"github.com/tomkukral/cloud-cleaner/registry"
)

func main() {
	cleanRegistry := false
	cleanCloudformation := true

	if cleanRegistry {
		url := "https://registry-1.docker.io/"
		username := ""
		password := ""
		repo := "kqueen/api"
		km := make([]string, 3)
		km[0] = "^v[0-9]+.[0-9]+$"
		km[1] = "^master$"
		km[2] = "^latest$"

		registry.CleanTags(url, repo, km, username, password)

	}

	if cleanCloudformation {
		cloudformation.CleanStacks(true)
	}

}
