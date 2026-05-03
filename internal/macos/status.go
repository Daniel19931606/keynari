//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit

#import <AppKit/AppKit.h>

@interface KeynariMenuTarget : NSObject
- (void)quit:(id)sender;
@end

@implementation KeynariMenuTarget
- (void)quit:(id)sender {
    [NSApp terminate:nil];
}
@end

static NSStatusItem *keynariStatusItem = nil;
static KeynariMenuTarget *keynariMenuTarget = nil;

static void runStatusApp(void) {
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
        NSMenuItem *status = [[NSMenuItem alloc] initWithTitle:@"Keynari is running" action:nil keyEquivalent:@""];
        [status setEnabled:NO];
        [menu addItem:status];
        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *quit = [[NSMenuItem alloc] initWithTitle:@"Quit Keynari" action:@selector(quit:) keyEquivalent:@"q"];
        [quit setTarget:keynariMenuTarget];
        [menu addItem:quit];

        [keynariStatusItem setMenu:menu];
        [NSApp run];
    }
}
*/
import "C"

// RunStatusApp blocks while the macOS menu bar app is running.
func RunStatusApp() {
	C.runStatusApp()
}
