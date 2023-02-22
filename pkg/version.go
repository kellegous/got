package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

var versionPat = regexp.MustCompile(`^(go)?(\d+)\.(\d+)\.(\d+)$`)

func NormalizeVersion(v string) (string, error) {
	if !versionPat.MatchString(v) {
		return "", fmt.Errorf("invalid version: %s", v)
	}

	if !strings.HasPrefix(v, "go") {
		return "go" + v, nil
	}

	return v, nil
}
