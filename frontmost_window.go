package main

import (
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"reflect"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
GetFrontMostAppPid(void){
	NSRunningApplication* app = [[NSWorkspace sharedWorkspace]
								  frontmostApplication];
	pid_t pid = [app processIdentifier];
	return pid;
}

CFStringRef
GetAppTitle(pid_t pid) {
	CFStringRef title = NULL;
	// Get the process ID of the frontmost application.
	//	NSRunningApplication* app = [[NSWorkspace sharedWorkspace]
	//								  frontmostApplication];
	//	pid_t pid = [app processIdentifier];

	// See if we have accessibility permissions, and if not, prompt the user to
	// visit System Preferences.
	NSDictionary *options = @{(id)kAXTrustedCheckOptionPrompt: @YES};
	Boolean appHasPermission = AXIsProcessTrustedWithOptions(
								 (__bridge CFDictionaryRef)options);
	if (!appHasPermission) {
	   return title; // we don't have accessibility permissions
	}
	// Get the accessibility element corresponding to the frontmost application.
	AXUIElementRef appElem = AXUIElementCreateApplication(pid);
	if (!appElem) {
	  return title;
	}

	// Get the accessibility element corresponding to the frontmost window
	// of the frontmost application.
	AXUIElementRef window = NULL;
	if (AXUIElementCopyAttributeValue(appElem,
		  kAXFocusedWindowAttribute, (CFTypeRef*)&window) != kAXErrorSuccess) {
	  CFRelease(appElem);
	  return title;
	}

	// Finally, get the title of the frontmost window.
	AXError result = AXUIElementCopyAttributeValue(window, kAXTitleAttribute,
					   (CFTypeRef*)&title);

	// At this point, we don't need window and appElem anymore.
	CFRelease(window);
	CFRelease(appElem);

	if (result != kAXErrorSuccess) {
	  // Failed to get the window title.
	  return title;
	}

	// Success! Now, do something with the title, e.g. copy it somewhere.

	// Once we're done with the title, release it.
	CFRelease(title);

    return title;
}
static inline CFIndex cfstring_utf8_length(CFStringRef str, CFIndex *need) {
  CFIndex n, usedBufLen;
  CFRange rng = CFRangeMake(0, CFStringGetLength(str));
  return CFStringGetBytes(str, rng, kCFStringEncodingUTF8, 0, 0, NULL, 0, need);
}
*/
import "C"

//import "github.com/shirou/gopsutil/v3/process"
func cfstringGo(cfs C.CFStringRef) (string, error) {
	if cfs == 0 {
		return "", errors.New("Null")
	}
	var usedBufLen C.CFIndex
	n := C.cfstring_utf8_length(cfs, &usedBufLen)
	if n <= 0 {
		return "", errors.New("null")
	}
	rng := C.CFRange{location: C.CFIndex(0), length: n}
	buf := make([]byte, int(usedBufLen))

	bufp := unsafe.Pointer(&buf[0])
	C.CFStringGetBytes(cfs, rng, C.kCFStringEncodingUTF8, 0, 0, (*C.UInt8)(bufp), C.CFIndex(len(buf)), &usedBufLen)

	sh := &reflect.StringHeader{
		Data: uintptr(bufp),
		Len:  int(usedBufLen),
	}
	return *(*string)(unsafe.Pointer(sh)), nil
}

type WindowInfo struct {
	Name       string
	Title      string
	Pid        int32
	Ppid       int32
	CreateTime int64
}

func GetActiveProcess() (map[string]interface{}, error) {
	dat := make(map[string]interface{})

	pid := C.GetFrontMostAppPid()
	if pid == 0 {
		return dat, errors.New("No")
	}
	dat["Pid"] = int32(pid)

	ps, _ := process.NewProcess(int32(pid))
	title_ref := C.CFStringRef(C.GetAppTitle(pid))
	title, err := cfstringGo(title_ref)
	if err != nil {
		return dat, err
	}
	dat["Title"] = title
	name, err := ps.Name()
	if err != nil {
		return dat, err
	}
	dat["Name"] = name

	create_time, err := ps.CreateTime()
	if err != nil {
		return dat, err
	}
	dat["CreateTime"] = create_time

	ppid, err := ps.Ppid()
	if err != nil {
		return dat, err
	}
	dat["Ppid"] = ppid

	return dat, nil
}

func main() {
	wi, err := GetActiveProcess()
	if err != nil {
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(wi)
	//fmt.Print(wi_str)
}
