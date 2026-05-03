#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP="$ROOT/dist/Keynari.app"
MACOS="$APP/Contents/MacOS"
RESOURCES="$APP/Contents/Resources"
ICON="$ROOT/dist/Keynari.icns"
BUILD="$ROOT/dist/build"

rm -rf "$APP"
mkdir -p "$MACOS" "$RESOURCES" "$BUILD"

go build -o "$MACOS/keynari-bin" "$ROOT/cmd/keynari"

if command -v magick >/dev/null 2>&1 && command -v iconutil >/dev/null 2>&1; then
	"$ROOT/scripts/make_icon.sh" >/dev/null
	cp "$ICON" "$RESOURCES/Keynari.icns"
fi

cat > "$BUILD/KeynariLauncher.swift" <<'SWIFT'
import AppKit
import UserNotifications

final class AppDelegate: NSObject, NSApplicationDelegate {
    private var statusItem: NSStatusItem!
    private var process: Process?

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
        startKeynari()
        buildMenu()
        notifyStarted()
    }

    func applicationWillTerminate(_ notification: Notification) {
        stopKeynari()
    }

    private func startKeynari() {
        let executable = Bundle.main.bundleURL
            .appendingPathComponent("Contents/MacOS/keynari-bin")
        let logFile = FileManager.default.homeDirectoryForCurrentUser
            .appendingPathComponent("Library/Logs/Keynari.log")

        let task = Process()
        task.executableURL = executable
        task.arguments = ["run", "--quiet", "--log-file", logFile.path]
        task.standardOutput = FileHandle.nullDevice
        task.standardError = FileHandle.nullDevice

        do {
            try task.run()
            process = task
        } catch {
            showAlert("Keynari could not start", error.localizedDescription)
        }
    }

    private func stopKeynari() {
        guard let process else { return }
        if process.isRunning {
            process.terminate()
            process.waitUntilExit()
        }
    }

    private func buildMenu() {
        statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
        if let button = statusItem.button {
            button.image = NSImage(systemSymbolName: "keyboard", accessibilityDescription: "Keynari")
            button.title = " Keynari"
        }

        let menu = NSMenu()
        let status = NSMenuItem(title: "Keynari is running", action: nil, keyEquivalent: "")
        status.isEnabled = false
        menu.addItem(status)
        menu.addItem(NSMenuItem.separator())
        menu.addItem(NSMenuItem(title: "Quit Keynari", action: #selector(quit), keyEquivalent: "q"))
        statusItem.menu = menu
    }

    private func notifyStarted() {
        UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .sound]) { granted, _ in
            guard granted else { return }
            let content = UNMutableNotificationContent()
            content.title = "Keynari is running"
            content.body = "Use the Keynari menu bar icon to quit."
            let request = UNNotificationRequest(identifier: UUID().uuidString, content: content, trigger: nil)
            UNUserNotificationCenter.current().add(request)
        }
    }

    private func showAlert(_ title: String, _ message: String) {
        let alert = NSAlert()
        alert.messageText = title
        alert.informativeText = message
        alert.alertStyle = .critical
        alert.runModal()
    }

    @objc private func quit() {
        NSApp.terminate(nil)
    }
}

let app = NSApplication.shared
let delegate = AppDelegate()
app.delegate = delegate
app.run()
SWIFT

xcrun swiftc "$BUILD/KeynariLauncher.swift" -o "$MACOS/Keynari" -framework AppKit -framework UserNotifications

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
