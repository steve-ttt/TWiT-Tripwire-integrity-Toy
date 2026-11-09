package main

import (
	"fmt"
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

func TestIntegrationLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	// Setup monitored directory + temp SQLite path + config JSON
	monitoredDir, _ := setupTestTree(t)
	dbPath := filepath.Join(t.TempDir(), "tripwire.db")
	rawConfig := fmt.Sprintf(`{
		"database": %q,
		"scan_dirs": [%q],
		"scan_files": []
	}`, dbPath, monitoredDir)

	cfg, err := readConfig([]byte(rawConfig))
	if err != nil {
		t.Fatalf("readConfig failed: %v", err)
	}

	// Initial "--update" (create baseline DB)
	updateBaseline(cfg)

	if failed := checkBaseline(cfg); failed {
		t.Fatal("expected clean integrity check to pass, but it failed")
	}

	// Tamper with a monitored file on disk
	tamperedFile := filepath.Join(monitoredDir, "file1.txt")
	if err := os.WriteFile(tamperedFile, []byte("hacked!"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Second "--check" should detect tampering (return true)
	if failed := checkBaseline(cfg); !failed {
		t.Fatal("expected integrity check to fail after tampering, but it passed")
	}
}
