//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework AppKit

#include <CoreGraphics/CoreGraphics.h>
#include <AppKit/AppKit.h>

void sendKeyEvent(CGKeyCode keyCode, bool keyDown) {
    CGEventRef event = CGEventCreateKeyboardEvent(NULL, keyCode, keyDown);
    CGEventPost(kCGHIDEventTap, event);
    CFRelease(event);
}

void sendKeyEventWithFlags(CGKeyCode keyCode, CGEventFlags flags, bool keyDown) {
    CGEventRef event = CGEventCreateKeyboardEvent(NULL, keyCode, keyDown);
    CGEventSetFlags(event, flags);
    CGEventPost(kCGHIDEventTap, event);
    CFRelease(event);
}

void setPasteboardString(char *text) {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSString *string = [NSString stringWithUTF8String:text];
        [pasteboard clearContents];
        [pasteboard setString:string forType:NSPasteboardTypeString];
    }
}

char *getPasteboardString(void) {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSString *string = [pasteboard stringForType:NSPasteboardTypeString];
        if (string == nil) {
            return NULL;
        }
        return strdup([string UTF8String]);
    }
}
*/
import "C"

import (
	"sync"
	"time"
	"unsafe"
)

const (
	keyDelete   = 0x33
	keyANSIv    = 0x09
	commandFlag = 0x00100000
)

type Replacer struct {
	mu sync.Mutex
}

func NewReplacer() *Replacer {
	return &Replacer{}
}

func (r *Replacer) Replace(oldLen int, newText string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	previousClipboard := C.getPasteboardString()
	if previousClipboard != nil {
		defer C.free(unsafe.Pointer(previousClipboard))
	}

	for i := 0; i < oldLen; i++ {
		C.sendKeyEvent(C.CGKeyCode(keyDelete), C.bool(true))
		C.sendKeyEvent(C.CGKeyCode(keyDelete), C.bool(false))
		time.Sleep(time.Millisecond)
	}

	cText := C.CString(newText)
	defer C.free(unsafe.Pointer(cText))
	C.setPasteboardString(cText)
	time.Sleep(15 * time.Millisecond)
	C.sendKeyEventWithFlags(C.CGKeyCode(keyANSIv), C.CGEventFlags(commandFlag), C.bool(true))
	C.sendKeyEventWithFlags(C.CGKeyCode(keyANSIv), C.CGEventFlags(commandFlag), C.bool(false))
	time.Sleep(120 * time.Millisecond)

	if previousClipboard != nil {
		C.setPasteboardString(previousClipboard)
	}

	return nil
}
