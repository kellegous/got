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

func requestFor(
	ctx context.Context,
	version string,
	platform *Platform,
) (*http.Request, error) {
	return http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://go.dev/gl/%s.%s-%s.tar.gz", version, platform.OS, platform.Arch),
		nil)
}

func Download(
	ctx context.Context,
	gotDir string,
	platform *Platform,
	version string,
	opts *DownloadOptions,
) error {
	name := fmt.Sprintf("%s-%s-%s", version, platform.OS, platform.Arch)
	dir := filepath.Join(gotDir, name)
	if _, err := os.Stat(dir); err == nil {
		if !opts.Force {
			return nil
		}

		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	if err := os.Mkdir(name, 0755); err != nil {
		return err
	}

	url := fmt.Sprintf("https://go.dev/gl/%s.%s-%s.tar.gz", version, platform.OS, platform.Arch)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d: %s", res.StatusCode, url)
	}

	return Untar(dir, res.Body)
}
