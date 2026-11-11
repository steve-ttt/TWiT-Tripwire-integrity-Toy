#!/usr/bin/env bash
set -e

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

DEMO_DIR="$TMP_DIR/demo"
CONFIG="$TMP_DIR/config.json"
DB="$TMP_DIR/tripwire.db"

mkdir -p "$DEMO_DIR"
echo "hello world" > "$DEMO_DIR/secret.txt"

cat << EOF > "$CONFIG"
{
    "database": "$DB",
    "scan_dirs": ["$DEMO_DIR"],
    "scan_files": []
}
EOF

echo "--- 1. Initializing baseline (--update) ---"
./bin/tripwire --update --config "$CONFIG"

echo "--- 2. Running clean check (--check) ---"
./bin/tripwire --check --config "$CONFIG"

echo "--- 3. Tampering with file and re-checking ---"
echo "tampered!" >> "$DEMO_DIR/secret.txt"
if ./bin/tripwire --check --config "$CONFIG"; then
    echo "ERROR: Expected integrity check to fail with non-zero exit code"
    exit 1
fi
echo "Smoke test passed!"
