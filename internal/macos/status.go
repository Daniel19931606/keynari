//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit

#import <AppKit/AppKit.h>
#include <stdlib.h>

@interface KeynariMenuTarget : NSObject
- (void)quit:(id)sender;
- (void)openAccessibility:(id)sender;
@end

@implementation KeynariMenuTarget
- (void)quit:(id)sender {
    [NSApp terminate:nil];
}
- (void)openAccessibility:(id)sender {
    NSURL *url = [NSURL URLWithString:@"x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"];
    [[NSWorkspace sharedWorkspace] openURL:url];
}
@end

static NSStatusItem *keynariStatusItem = nil;
static KeynariMenuTarget *keynariMenuTarget = nil;
static NSMenuItem *keynariStatusMenuItem = nil;

static void setStatusTitle(const char *title) {
    if (keynariStatusMenuItem == nil) {
        return;
    }
    NSString *value = [NSString stringWithUTF8String:title];
    dispatch_async(dispatch_get_main_queue(), ^{
        [keynariStatusMenuItem setTitle:value];
    });
}

static void runStatusApp(const char *initialStatus) {
    @autoreleasepool {
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

        keynariMenuTarget = [KeynariMenuTarget new];
        keynariStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];

        NSStatusBarButton *button = [keynariStatusItem button];
        [button setTitle:@" Keynari"];

        if (@available(macOS 11.0, *)) {
            NSImage *image = [NSImage imageWithSystemSymbolName:@"keyboard" accessibilityDescription:@"Keynari"];
            [button setImage:image];
        }

        NSMenu *menu = [NSMenu new];
        NSString *statusTitle = [NSString stringWithUTF8String:initialStatus];
        keynariStatusMenuItem = [[NSMenuItem alloc] initWithTitle:statusTitle action:nil keyEquivalent:@""];
        [keynariStatusMenuItem setEnabled:NO];
        [menu addItem:keynariStatusMenuItem];
        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *openAccessibility = [[NSMenuItem alloc] initWithTitle:@"Open Accessibility Settings" action:@selector(openAccessibility:) keyEquivalent:@""];
        [openAccessibility setTarget:keynariMenuTarget];
        [menu addItem:openAccessibility];

        NSMenuItem *quit = [[NSMenuItem alloc] initWithTitle:@"Quit Keynari" action:@selector(quit:) keyEquivalent:@"q"];
        [quit setTarget:keynariMenuTarget];
        [menu addItem:quit];

        [keynariStatusItem setMenu:menu];
        [NSApp run];
    }
}
*/
import "C"
import "unsafe"

// RunStatusApp blocks while the macOS menu bar app is running.
func RunStatusApp(initialStatus string) {
	cStatus := C.CString(initialStatus)
	defer C.free(unsafe.Pointer(cStatus))
	C.runStatusApp(cStatus)
}

func SetStatusTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.setStatusTitle(cTitle)
}
