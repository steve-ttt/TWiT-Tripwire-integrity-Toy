package tripwire

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanDirsAndFiles(t *testing.T) {
	dir, expectedFiles := setupTestTree(t)
	dbPath := filepath.Join(t.TempDir(), "test.db")
	saveBaseline(dbPath, expectedFiles)

	scanResult, err := hashTree(dir)
	if err != nil {
		t.Fatalf("failed to scan directory: %v", err)
	}
	storedFiles, err := loadBaseline(dbPath)
	if err != nil {
		t.Fatalf("failed to load baseline: %v", err)
	}

	if !reflect.DeepEqual(scanResult, storedFiles) {
		t.Errorf("scan result does not match stored files")
	}
}

func TestScanDirsAndFilesWithUnknowFile(t *testing.T) {
	dir, expectedFiles := setupTestTree(t)
	// Add an unknown file to the scan result
	unknownFilePath := filepath.Join(dir, "unknown.txt")
	if err := os.WriteFile(unknownFilePath, []byte("unknown content"), 0o600); err != nil {
		t.Fatalf("failed to create unknown file: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "test.db")
	saveBaseline(dbPath, expectedFiles)

	scanResult, err := hashTree(dir)
	if err != nil {
		t.Fatalf("failed to scan directory: %v", err)
	}
	storedFiles, err := loadBaseline(dbPath)
	if err != nil {
		t.Fatalf("failed to load baseline: %v", err)
	}

	if reflect.DeepEqual(scanResult, storedFiles) {
		t.Errorf("scan result should not match stored files due to unknown file")
	}
}
