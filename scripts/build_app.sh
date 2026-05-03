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
import Darwin

final class AppDelegate: NSObject, NSApplicationDelegate {
    private var statusItem: NSStatusItem!
    private var process: Process?
    private var statusMenuItem: NSMenuItem!
    private var restartMenuItem: NSMenuItem!
    private var launchAttemptAt: Date?

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
        buildMenu()
        startKeynari()
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
        task.terminationHandler = { [weak self] finishedTask in
            DispatchQueue.main.async {
                guard let self else { return }
                if self.process === finishedTask {
                    let livedFor = self.launchAttemptAt.map { Date().timeIntervalSince($0) } ?? 0
                    self.process = nil
                    self.launchAttemptAt = nil
                    if livedFor < 1.5 {
                        self.setStatus("Keynari needs Accessibility permission")
                    } else {
                        self.setStatus("Keynari is stopped")
                    }
                }
            }
        }

        do {
            launchAttemptAt = Date()
            try task.run()
            process = task
            setStatus("Keynari is running")
        } catch {
            setStatus("Keynari failed to start")
            showAlert("Keynari could not start", error.localizedDescription)
        }
    }

    private func stopKeynari() {
        guard let process else { return }
        defer {
            self.process = nil
            self.launchAttemptAt = nil
        }

        guard process.isRunning else { return }

        let pid = process.processIdentifier
        process.interrupt()

        for _ in 0..<15 {
            if !process.isRunning {
                return
            }
            usleep(100_000)
        }

        if kill(pid, 0) == 0 {
            _ = kill(pid, SIGTERM)
        }

        for _ in 0..<10 {
            if !process.isRunning {
                return
            }
            usleep(100_000)
        }

        if kill(pid, 0) == 0 {
            _ = kill(pid, SIGKILL)
        }

        let killCommand = """
        pkill -TERM -f 'Contents/MacOS/keynari-bin run --quiet --log-file'
        sleep 0.3
        pkill -KILL -f 'Contents/MacOS/keynari-bin run --quiet --log-file' || true
        """
        let cleanup = Process()
        cleanup.executableURL = URL(fileURLWithPath: "/bin/sh")
        cleanup.arguments = ["-c", killCommand]
        try? cleanup.run()
        cleanup.waitUntilExit()
    }

    private func buildMenu() {
        statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
        if let button = statusItem.button {
            button.image = NSImage(systemSymbolName: "keyboard", accessibilityDescription: "Keynari")
            button.title = " Keynari"
        }

        let menu = NSMenu()
        statusMenuItem = NSMenuItem(title: "Keynari is starting...", action: nil, keyEquivalent: "")
        statusMenuItem.isEnabled = false
        menu.addItem(statusMenuItem)
        menu.addItem(NSMenuItem.separator())
        restartMenuItem = NSMenuItem(title: "Restart Keynari", action: #selector(restart), keyEquivalent: "r")
        menu.addItem(restartMenuItem)
        menu.addItem(NSMenuItem(title: "Quit Keynari", action: #selector(quit), keyEquivalent: "q"))
        statusItem.menu = menu
    }

    private func setStatus(_ title: String) {
        statusMenuItem?.title = title
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
        stopKeynari()
        NSApp.terminate(nil)
    }

    @objc private func restart() {
        stopKeynari()
        startKeynari()
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
