#include "ns_window.h"

void Center(void *self) {
    NSWindow *window = self;
    [window center];
}

void SetTitle(void *self, char *title) {
    NSWindow *window = self;
    NSString *nsTitle =  [NSString stringWithUTF8String:title];

    [window setTitle:nsTitle];
    free(title);
}