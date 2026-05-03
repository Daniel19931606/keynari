//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ApplicationServices

#include <CoreGraphics/CoreGraphics.h>
#include <ApplicationServices/ApplicationServices.h>

extern void goKeyEventCallback(int64_t keycode, uint64_t flags, UniChar character);

static CFMachPortRef globalTap = NULL;

static CGEventRef eventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        if (globalTap != NULL) {
            CGEventTapEnable(globalTap, true);
        }
        return event;
    }

    if (type == kCGEventKeyDown) {
        CGKeyCode keycode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
        CGEventFlags flags = CGEventGetFlags(event);
        UniChar chars[4];
        UniCharCount actualLength = 0;
        CGEventKeyboardGetUnicodeString(event, 4, &actualLength, chars);

        UniChar character = 0;
        if (actualLength > 0) {
            character = chars[0];
        }

        goKeyEventCallback((int64_t)keycode, (uint64_t)flags, character);
    }

    return event;
}

static bool checkAccessibilityPermissions(void) {
    return AXIsProcessTrusted();
}

static void requestAccessibilityPermissions(void) {
    NSDictionary *options = @{(__bridge id)kAXTrustedCheckOptionPrompt: @YES};
    AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
}

static CFMachPortRef createEventTap(void) {
    CGEventMask eventMask = (1 << kCGEventKeyDown);
    globalTap = CGEventTapCreate(
        kCGSessionEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionListenOnly,
        eventMask,
        eventCallback,
        NULL
    );
    return globalTap;
}

static void enableEventTap(CFMachPortRef tap, int enable) {
    CGEventTapEnable(tap, enable ? true : false);
}

static bool isMachPortNull(CFMachPortRef ref) {
    return ref == NULL;
}

static bool isRunLoopSourceNull(CFRunLoopSourceRef ref) {
    return ref == NULL;
}

static bool isRunLoopNull(CFRunLoopRef ref) {
    return ref == NULL;
}

static CFMachPortRef nullMachPort(void) {
    return NULL;
}

static CFRunLoopSourceRef nullRunLoopSource(void) {
    return NULL;
}
*/
import "C"

import (
	"errors"
	"sync"
	"time"
)

const (
	KeyBackspace = 0x33
)

type KeyEvent struct {
	Char      rune
	KeyCode   uint16
	Modifiers Modifiers
	Timestamp time.Time
}

type Modifiers struct {
	Shift bool
	Ctrl  bool
	Alt   bool
	Meta  bool
}

type Listener struct {
	events     chan KeyEvent
	eventTap   C.CFMachPortRef
	runLoopSrc C.CFRunLoopSourceRef
	runLoop    C.CFRunLoopRef
	mu         sync.Mutex
}

var (
	activeListener *Listener
	listenerMu     sync.Mutex
)

func NewListener() *Listener {
	return &Listener{events: make(chan KeyEvent, 512)}
}

func EnsureAccessibility() error {
	if bool(C.checkAccessibilityPermissions()) {
		return nil
	}
	return errors.New("accessibility permission required: grant access to Keynari in System Settings > Privacy & Security > Accessibility, and if you rebuilt the app remove the old entry and add the new one again")
}

func RequestAccessibility() {
	C.requestAccessibilityPermissions()
}

func (l *Listener) Events() <-chan KeyEvent {
	return l.events
}

func (l *Listener) Start() error {
	if err := EnsureAccessibility(); err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	listenerMu.Lock()
	activeListener = l
	listenerMu.Unlock()

	l.eventTap = C.createEventTap()
	if bool(C.isMachPortNull(l.eventTap)) {
		return errors.New("failed to create keyboard event tap")
	}

	l.runLoopSrc = C.CFMachPortCreateRunLoopSource(C.kCFAllocatorDefault, l.eventTap, 0)
	if bool(C.isRunLoopSourceNull(l.runLoopSrc)) {
		return errors.New("failed to create run-loop source")
	}

	go l.run()
	return nil
}

func (l *Listener) run() {
	l.runLoop = C.CFRunLoopGetCurrent()
	C.CFRunLoopAddSource(l.runLoop, l.runLoopSrc, C.kCFRunLoopCommonModes)
	C.enableEventTap(l.eventTap, 1)
	C.CFRunLoopRun()
}

func (l *Listener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !bool(C.isRunLoopNull(l.runLoop)) {
		C.CFRunLoopStop(l.runLoop)
	}
	if !bool(C.isMachPortNull(l.eventTap)) {
		C.enableEventTap(l.eventTap, 0)
		C.CFRelease(C.CFTypeRef(l.eventTap))
		l.eventTap = C.nullMachPort()
	}
	if !bool(C.isRunLoopSourceNull(l.runLoopSrc)) {
		C.CFRelease(C.CFTypeRef(l.runLoopSrc))
		l.runLoopSrc = C.nullRunLoopSource()
	}

	listenerMu.Lock()
	if activeListener == l {
		activeListener = nil
	}
	listenerMu.Unlock()
	close(l.events)
}

//export goKeyEventCallback
func goKeyEventCallback(keycode C.int64_t, flags C.uint64_t, character C.UniChar) {
	listenerMu.Lock()
	l := activeListener
	listenerMu.Unlock()
	if l == nil {
		return
	}

	const (
		shiftFlag   = 0x00020000
		controlFlag = 0x00040000
		altFlag     = 0x00080000
		commandFlag = 0x00100000
	)

	event := KeyEvent{
		Char:    rune(character),
		KeyCode: uint16(keycode),
		Modifiers: Modifiers{
			Shift: uint64(flags)&shiftFlag != 0,
			Ctrl:  uint64(flags)&controlFlag != 0,
			Alt:   uint64(flags)&altFlag != 0,
			Meta:  uint64(flags)&commandFlag != 0,
		},
		Timestamp: time.Now(),
	}

	select {
	case l.events <- event:
	default:
	}
}
