package pkg

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadOptions struct {
	Force bool
}

func Download(
	ctx context.Context,
	gotDir string,
	platform *Platform,
	version string,
	opts *DownloadOptions,
) (string, error) {
	name := fmt.Sprintf("%s-%s-%s", version, platform.OS, platform.Arch)
	dir := filepath.Join(gotDir, name)
	if _, err := os.Stat(dir); err == nil {
		if !opts.Force {
			return name, nil
		}

		if err := os.RemoveAll(dir); err != nil {
			return "", err
		}
	}

	if err := os.Mkdir(dir, 0755); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://go.dev/dl/%s.%s-%s.tar.gz", version, platform.OS, platform.Arch)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
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
		return "", fmt.Errorf("http status %d: %s", res.StatusCode, url)
	}

	if err := Untar(dir, res.Body); err != nil {
		return "", err
	}

	return name, nil
}
