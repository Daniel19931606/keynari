#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${1:-dev}"
RELEASE_DIR="$ROOT/dist/release"
APP="$ROOT/dist/Keynari.app"
ZIP="$RELEASE_DIR/Keynari-macOS-${VERSION}.zip"
CHECKSUMS="$RELEASE_DIR/checksums.txt"

"$ROOT/scripts/build_app.sh"

rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

ditto -c -k --keepParent "$APP" "$ZIP"

(
	cd "$RELEASE_DIR"
	shasum -a 256 "$(basename "$ZIP")" > "$CHECKSUMS"
)

echo "Packaged $ZIP"
echo "Checksums $CHECKSUMS"
