#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP="$ROOT/dist/Keynari.app"
MACOS="$APP/Contents/MacOS"
RESOURCES="$APP/Contents/Resources"
ICON="$ROOT/dist/Keynari.icns"

rm -rf "$APP"
mkdir -p "$MACOS" "$RESOURCES"

go build -o "$MACOS/Keynari" "$ROOT/cmd/keynari"

if command -v magick >/dev/null 2>&1 && command -v iconutil >/dev/null 2>&1; then
	"$ROOT/scripts/make_icon.sh" >/dev/null
	cp "$ICON" "$RESOURCES/Keynari.icns"
fi
cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "https://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>en</string>
	<key>CFBundleDisplayName</key>
	<string>Keynari</string>
	<key>CFBundleExecutable</key>
	<string>Keynari</string>
	<key>CFBundleIdentifier</key>
	<string>com.daniel19931606.keynari</string>
	<key>CFBundleIconFile</key>
	<string>Keynari</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>Keynari</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>0.1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSMinimumSystemVersion</key>
	<string>12.0</string>
	<key>LSUIElement</key>
	<true/>
	<key>NSAppleEventsUsageDescription</key>
	<string>Keynari needs accessibility access to replace mistyped words in the active app.</string>
	<key>NSUserNotificationAlertStyle</key>
	<string>alert</string>
</dict>
</plist>
PLIST

if command -v codesign >/dev/null 2>&1; then
	codesign --force --deep --sign - "$APP" >/dev/null
fi

echo "Built $APP"
