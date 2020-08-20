package cocoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
//#include "ns_window.h"
import "C"
import "unsafe"

type NSWindow struct {
	Ptr unsafe.Pointer
}

func (self *NSWindow) Center() {
	C.Center(self.Ptr)
}

func (self *NSWindow) SetTitle(title string) {
	C.SetTitle(self.Ptr, C.CString(title))
}
