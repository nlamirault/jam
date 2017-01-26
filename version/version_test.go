package version

import (
	// "fmt"
	"regexp"
	"testing"
)

func Test_Version(t *testing.T) {
	validVersion := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	if !validVersion.MatchString(Version) {
		t.Fatal("Invalid version")
	}
}
