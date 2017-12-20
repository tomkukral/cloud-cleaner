package registry

import (
	"fmt"
	"testing"
)

func TestCheckMatch(t *testing.T) {
	matches := []string{
		"^master$",
		"^latest$",
		"^v[0-9]+.[0-9]+$",
	}

	names := []struct {
		name string
		req  bool
	}{
		{"master", true},
		{"latest", true},
		{"v0.1", true},
		{"latestlatest", false},
		{"v1.1.1", false},
	}

	for _, tt := range names {
		fmt.Println(tt)
		res := checkMatch(matches, tt.name)
		if res != tt.req {
			t.Errorf("checkMatch for name '%s': expected '%t', got '%t'", tt.name, tt.req, res)
		}
	}

}
