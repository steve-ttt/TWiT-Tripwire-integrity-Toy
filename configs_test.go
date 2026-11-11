package main

import (
	"testing"
)

func TestConfigStruct(t *testing.T) {
	in := []byte(`{
		"database": "test.db",
		"scan_dirs": ["/path/to/dir1", "/path/to/dir2"],
		"scan_files": ["/path/to/file1", "/path/to/file2"]
	}`)

	got, err := readConfig(in)
	if err != nil {
		t.Fatalf("readConfig returned an error: %v", err)
	}

	if got.Database != "test.db" {
		t.Errorf("Database = %q, want %q", got.Database, "test.db")
	}
	if len(got.ScanDirs) != 2 || got.ScanDirs[0] != "/path/to/dir1" || got.ScanDirs[1] != "/path/to/dir2" {
		t.Errorf("ScanDirs = %v, want [%q, %q]", got.ScanDirs, "/path/to/dir1", "/path/to/dir2")
	}
	if len(got.ScanFiles) != 2 || got.ScanFiles[0] != "/path/to/file1" || got.ScanFiles[1] != "/path/to/file2" {
		t.Errorf("ScanFiles = %v, want [%q, %q]", got.ScanFiles, "/path/to/file1", "/path/to/file2")
	}

}

func TestReadConfigInvalidJSON(t *testing.T) {
	_, err := readConfig([]byte(`{invalid json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
