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
	"context"
	"runtime"
	"sync"
	"unsafe"
)

/*********************MobileDevice********************/
type NotificationHander func(even, mode string, modeDevice, restorableDevice uintptr)

var (
	mutex            sync.Mutex
	registedHandle   []NotificationHander
	loopCancleForWin context.CancelFunc
	mainLoopThread   uintptr
)

//AMRestorableDeviceRegisterForNotificationsForDevices ...
func AMRestorableDeviceRegisterForNotificationsForDevices(hander NotificationHander) {
	if hander == nil {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	if registedHandle == nil {
		registedHandle = make([]NotificationHander, 0, 5)
	}
	registedHandle = append(registedHandle, hander)
	C.AMRestorableDeviceRegisterForNotificationsForDevices((*[0]byte)(unsafe.Pointer(C.restoreDeviceNotify)), nil, 0x4F, nil, nil)
}

//export restoreDeviceNotify
func restoreDeviceNotify(device *C.struct_am_restore_device, even C.int) {
	if nil == device {
		return
	}

	// 4 -- normal mode
	// 2 -- recovery mode
	// 1 -- dfu mode
	state := C.AMRestorableDeviceGetState(unsafe.Pointer(device))
	mode := ""
	ev := ""
	if even == 0 {
		ev = "insert"
	} else {
		ev = "extract"
	}

	var modeDevice unsafe.Pointer
	if state == 2 {
		modeDevice = C.AMRestorableDeviceCopyRecoveryModeDevice(unsafe.Pointer(device))
		mode = "recovery"
	} else if state == 1 {
		modeDevice = C.AMRestorableDeviceCopyDFUModeDevice(unsafe.Pointer(device))
		mode = "dfu"
	} else if state == 4 {
		modeDevice = C.AMRestorableDeviceCopyAMDevice(unsafe.Pointer(device))
		mode = "normal"
	}
	mutex.Lock()
	tmpHandler := registedHandle
	mutex.Unlock()
	for _, f := range tmpHandler {
		f(ev, mode, uintptr(modeDevice), uintptr(unsafe.Pointer(device)))
	}
}

//RunLoopRun ...
func RunLoopRun() {
	if runtime.GOOS != "windows" {
		runtime.LockOSThread()
		mainLoopThread = uintptr(C.CFRunLoopGetCurrent())
		C.CFRunLoopRun()
	} else {
		if loopCancleForWin != nil {
			return
		}
		var contextForWin context.Context
		contextForWin, loopCancleForWin = context.WithCancel(context.Background())
		select {
		case <-contextForWin.Done():
		}
	}
}

//RunLoopStop ...
func RunLoopStop() {
	if runtime.GOOS != "windows" {
		if 0 != mainLoopThread {
			C.CFRunLoopStop(C.CFRunLoopRef(mainLoopThread))
		}
	} else {
		if loopCancleForWin != nil {
			loopCancleForWin()
			loopCancleForWin = nil
		}
	}
}
