package tripwire

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashFile(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    string
	}{
		{
			name:    "empty",
			content: []byte(""),
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:    "abc",
			content: []byte("abc"),
			want:    "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
	}

	dir := t.TempDir() // automatically removed when TestHashFile finishes

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name)
			if err := os.WriteFile(path, tt.content, 0o600); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			got, err := hashFile(path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHashFileNotFound(t *testing.T) {
	_, err := hashFile("nonexistent_file.txt")
	if err == nil {
		t.Fatal("expected an error for a nonexistent file, got nil")
	}
}

func TestHashTree(t *testing.T) {
	dir := t.TempDir() // automatically removed when TestHashTree finishes

	// Create a simple directory structure
	subDir := filepath.Join(dir, "subdir")
	if err := os.Mkdir(subDir, 0o700); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	// Create files in the root directory
	file1Path := filepath.Join(dir, "file1.txt")
	if err := os.WriteFile(file1Path, []byte("file1 content"), 0o600); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	file2Path := filepath.Join(dir, "file2.txt")
	if err := os.WriteFile(file2Path, []byte("file2 content"), 0o600); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}
	// Create a file in the subdirectory
	subFilePath := filepath.Join(subDir, "subfile.txt")
	if err := os.WriteFile(subFilePath, []byte("subfile content"), 0o600); err != nil {
		t.Fatalf("failed to create subfile: %v", err)
	}

	hashes, err := hashTree(dir)
	if err != nil || hashes == nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the hashes map contains the expected files and their hashes
	expectedFiles := map[string]string{
		file1Path:   "1705789d380ee110bc09231df8af42a0cc564a1510ebd2168516d4985c40a263",
		file2Path:   "69d21bef329d53a0f62bf6f19faab99eab3dac4d53e43fbbb257fee4d14f60ab",
		subFilePath: "f53617a5043bff46ec5fac4874d9cd0fc8bd50b57620488597ed630a41976e58",
	}

	for path, hash := range hashes {
		expectedhash, ok := expectedFiles[path]
		if !ok {
			t.Errorf("unexpected file in hashes: %s", path)
			continue
		}
		if hash != expectedhash {
			t.Errorf("hash for %s: got %q, want %q", path, hash, expectedhash)
		}
	}

}
