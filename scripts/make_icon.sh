#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ICONSET="$ROOT/dist/Keynari.iconset"
ICON="$ROOT/dist/Keynari.icns"
BASE="$ROOT/dist/keynari-icon-1024.png"

rm -rf "$ICONSET" "$ICON" "$BASE"
mkdir -p "$ICONSET"

magick -size 1024x1024 gradient:"#36D1DC-#5B86E5" -rotate 35 -gravity center -crop 1024x1024+0+0 +repage \
	-fill "rgba(255,255,255,0.18)" -draw "circle 260,230 500,230" \
	-fill "rgba(0,0,0,0.18)" -draw "roundrectangle 168,168 856,856 196,196" \
	-fill "rgba(255,255,255,0.95)" -draw "roundrectangle 196,196 828,828 164,164" \
	-fill "#111827" -font "/System/Library/Fonts/Supplemental/Arial Bold.ttf" -pointsize 470 -gravity center -annotate +0-10 "K" \
	-fill "#F97316" -draw "roundrectangle 612,610 802,800 56,56" \
	-fill white -font "/System/Library/Fonts/Supplemental/Arial Bold.ttf" -pointsize 120 -gravity southeast -annotate +254+240 "Я" \
	"$BASE"

sips -z 16 16 "$BASE" --out "$ICONSET/icon_16x16.png" >/dev/null
sips -z 32 32 "$BASE" --out "$ICONSET/icon_16x16@2x.png" >/dev/null
sips -z 32 32 "$BASE" --out "$ICONSET/icon_32x32.png" >/dev/null
sips -z 64 64 "$BASE" --out "$ICONSET/icon_32x32@2x.png" >/dev/null
sips -z 128 128 "$BASE" --out "$ICONSET/icon_128x128.png" >/dev/null
sips -z 256 256 "$BASE" --out "$ICONSET/icon_128x128@2x.png" >/dev/null
sips -z 256 256 "$BASE" --out "$ICONSET/icon_256x256.png" >/dev/null
sips -z 512 512 "$BASE" --out "$ICONSET/icon_256x256@2x.png" >/dev/null
sips -z 512 512 "$BASE" --out "$ICONSET/icon_512x512.png" >/dev/null
cp "$BASE" "$ICONSET/icon_512x512@2x.png"

iconutil -c icns "$ICONSET" -o "$ICON"
echo "Built $ICON"
