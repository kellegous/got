package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var versionPat = regexp.MustCompile(`^(go)?(\d+)\.(\d+)\.(\d+)$`)

type Versions struct {
	latest string
}

func (v *Versions) Normalize(
	ctx context.Context,
	version string,
) (string, error) {
	var err error
	if version == "latest" {
		if v.latest == "" {
			v.latest, err = getLatestVersion(ctx)
			if err != nil {
				return "", err
			}
		}
		return v.latest, nil
	}

	if !versionPat.MatchString(version) {
		return "", fmt.Errorf("invalid version: %s", version)
	}

	if !strings.HasPrefix(version, "go") {
		return "go" + version, nil
	}

	return version, nil
}

func getLatestVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://go.dev/dl/?mode=json",
		nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d (%s)", res.StatusCode, res.Status)
	}

	var versions []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	if err := json.NewDecoder(res.Body).Decode(&versions); err != nil {
		return "", err
	}

	for _, version := range versions {
		if version.Stable {
			return version.Version, nil
		}
	}

	return "", errors.New("no stable version found")
}
