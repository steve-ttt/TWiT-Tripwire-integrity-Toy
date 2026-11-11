package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func hashTree(root string) (map[string]string, error) {
	hashes := make(map[string]string)
	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() { // ignore directories and non-regular files (ie. symlinks, sockets, etc.)
			hash, err := hashFile(path)
			if err != nil {
				return err
			}
			hashes[path] = hash
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return hashes, nil
}
