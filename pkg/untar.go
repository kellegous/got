package pkg

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func copyTo(dst string, r io.Reader, perm os.FileMode) error {
	w, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

func Untar(dst string, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	for tr := tar.NewReader(gr); ; {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		rel, err := filepath.Rel("go", header.Name)
		if err != nil {
			return err
		}

		path := filepath.Join(dst, rel)

		if m := header.FileInfo().Mode(); m.IsDir() {
			if _, err := os.Stat(path); err != nil {
				if err := os.MkdirAll(path, m.Perm()); err != nil {
					return err
				}
			}
		} else if m.IsRegular() {
			if err := copyTo(path, tr, m.Perm()); err != nil {
				return err
			}
		}
	}
}
