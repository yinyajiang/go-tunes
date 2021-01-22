package main

/*
#cgo windows CPPFLAGS: -DWIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation

#include "simpleApi.h"
extern void restoreDeviceNotify(struct am_restore_device* pDevice, int nMode);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	C.AMRestorableDeviceRegisterForNotificationsForDevices((*[0]byte)(unsafe.Pointer(C.restoreDeviceNotify)), nil, 0x4F, nil, nil)
	select {}
}

//export restoreDeviceNotify
func restoreDeviceNotify(device *C.struct_am_restore_device, even C.int) {
	if nil == device {
		return
	}

	// 4 -- normal mode
	// 2 -- recovery mode
	// 1 -- dfu mode
	mode := C.AMRestorableDeviceGetState(unsafe.Pointer(device))
	strmode := ""
	streven := ""
	if even == 0 {
		streven = "insert"
	} else {
		streven = "extract"
	}

	var modeDevice unsafe.Pointer
	if mode == 2 {
		modeDevice = C.AMRestorableDeviceCopyRecoveryModeDevice(unsafe.Pointer(device))
		strmode = "recovery"
	} else if mode == 1 {
		modeDevice = C.AMRestorableDeviceCopyDFUModeDevice(unsafe.Pointer(device))
		strmode = "dfu"
	} else if mode == 4 {
		modeDevice = C.AMRestorableDeviceCopyAMDevice(unsafe.Pointer(device))
		strmode = "normal"
	}
	fmt.Println(strmode, " ", streven, " ", modeDevice)
}
