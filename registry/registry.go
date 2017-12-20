package registry

import (
	"fmt"
	"os"
	"regexp"

	"github.com/heroku/docker-registry-client/registry"
)

func checkMatch(matches []string, name string) bool {
	for i := 0; i < len(matches); i++ {
		result, err := regexp.MatchString(matches[i], name)
		if result && err == nil {
			return true
		}
	}

	return false
}

func makeError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func CleanTags(url string, repo string, km []string, username string, password string) {

	hub, err := registry.New(url, username, password)

	if err != nil {
		makeError("Unable to connect to registry")
	}

	tags, err := hub.Tags(repo)
	if err != nil {
		makeError("Unable to get tags")
	}

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		fmt.Printf(tag)

		if !checkMatch(km, tag) {
			fmt.Printf(" DELETE\n")
			digest, err := hub.ManifestDigest(repo, tag)
			if err != nil {
				makeError("Unable to read manifest for tag")
			}

			err = hub.DeleteManifest(repo, digest)
			if err != nil {
				fmt.Println(err)
				makeError("Failed to delete tag")
			} else {
				fmt.Println("Tag %s:%s deleted", repo, tag)
			}
		}
	}
}
