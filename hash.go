package tripwire

import (
	"crypto/sha256"
	"fmt"
	"io"
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
	// TDD stub
	hashes := make(map[string]string)
	e := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() { // ignore directories and non-regular files (ie. symlinks, sockets, etc.)
			hash, err := hashFile(path)
			if err != nil {
				return err
			}
			hashes[path] = hash
			// fmt.Printf("File: %s, Hash: %s\n", path, hash)
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return hashes, nil
}
