package tripwire

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := createDatabase(dbPath)
	if err != nil {
		t.Fatalf("createDatabase failed: %v", err)
	}
	defer db.Close()

	// Check if the database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("Database file was not created at %s", dbPath)
	}

	// Check if the baseline table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='baseline';").Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to query for baseline table: %v", err)
	}
	if tableName != "baseline" {
		t.Fatalf("Expected table 'baseline', got '%s'", tableName)
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	want := map[string]string{
		"/etc/passwd": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"/etc/hosts":  "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
	}

	if err := saveBaseline(dbPath, want); err != nil {
		t.Fatalf("saveBaseline failed: %v", err)
	}

	got, err := loadBaseline(dbPath)
	if err != nil {
		t.Fatalf("loadBaseline failed: %v", err)
	}

	// compare got vs want...
	if len(got) != len(want) {
		t.Errorf("Expected %d entries, got %d", len(want), len(got))
	}
	for k, v := range want {
		if gotV, exists := got[k]; !exists {
			t.Errorf("Expected key %s not found", k)
		} else if gotV != v {
			t.Errorf("For key %s, expected %s, got %s", k, v, gotV)
		}
	}
}
