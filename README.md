# Trip Wire integrity Toy - TWiT 

## ! Warning this is a toy it is not suitable for production !

A TOY file integrity monitoring (FIM) tool written in Go. `tripWire` computes SHA-256 hashes of configured files and directory trees, stores a known-good baseline in a local SQLite database, and alerts on unauthorized file additions, modifications, or deletions.

## But Why ?

This is a toy implementation to learn how the unix/linux open source tool TripWire works under the hood. Also I wanted to practice GO-lang and TDD, because reasons 🤔 
Lastly, of course for fun !

## Features

- **SHA-256 File & Directory Hashing**: Recursively scans regular files in configured directories as well as individual standalone files.
- **Zero-CGO SQLite Baseline Storage**: Persists baseline file paths and SHA-256 digests in a local SQLite database using `modernc.org/sqlite`.
- **Tamper Detection**: Reports newly added unknown files, modified files (digest mismatch), and deleted files compared against the baseline.
- **Simple JSON Configuration**: Configure monitored directories, individual files, and the database path in a single JSON file.

## Getting Started

### Prerequisites

- Go 1.26+

### Build

Use the included `Makefile` or `go build` directly:

```bash
make build
```

This compiles the binary to `./bin/tripwire`.

## Configuration

Create a JSON configuration file (see `config.example.json`):

```json
{
    "database": "/root/tripwire.db",
    "scan_dirs": ["/etc", "/path/to/dir2"],
    "scan_files": ["/boot/grub2/grub.cfg", "/path/to/file2"]
}
```

| Field | Description |
| :--- | :--- |
| `database` | File path to the SQLite database used to store the baseline hashes. |
| `scan_dirs` | List of directory paths to recursively scan (regular files only). |
| `scan_files` | List of individual file paths to monitor. |

By default, `tripwire` looks for `/root/tripwire.json` unless overridden via `--config`.

## Usage

### 1. Initialize or Update the Baseline

Scan all configured directories and files and save their SHA-256 hashes to the database:

```bash
./bin/tripwire --update --config /path/to/config.json
```

*(Omit `--config <file>` to use the default `/root/tripwire.json`.)*

### 2. Check File Integrity

Compare the current filesystem state against the saved baseline:

```bash
./bin/tripwire --check --config /path/to/config.json
```

Example output when tampering is detected:

```text
File changed: /tmp/demo/secret.txt
==================================
Integrity check failed.
==================================
```

### 3. Help

```bash
./bin/tripwire --help
```

## Testing

Run unit tests, integration tests, and the end-to-end smoke test suite via `make`:

```bash
make test
```

Or run them individually:

```bash
# Run unit and integration tests
go test -v ./...

# Run unit tests only (skip integration lifecycle test)
go test -short ./...

# Run end-to-end CLI smoke test
./smoke_test.sh
```

## Project Structure

- `main.go` — CLI entrypoint, directory/file scanning orchestration, and baseline comparison (`--check` / `--update`).
- `configs.go` — JSON configuration struct (`Config`) and parser.
- `hash.go` — SHA-256 hashing for single files (`hashFile`) and recursive directory trees (`hashTree`).
- `database.go` — SQLite schema creation and baseline persistence (`createDatabase`, `saveBaseline`, `loadBaseline`).
- `smoke_test.sh` — End-to-end shell test demonstrating baseline creation, clean check, and tamper detection.
