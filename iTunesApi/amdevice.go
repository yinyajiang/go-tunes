package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

#include "simpleApi.h"
extern void restoreDeviceNotify(struct am_restore_device* pDevice, int nMode);
*/
import "C"
import (
	"context"
	"sync"
	"unsafe"
)

type NotificationHander func(even, mode string, modeDevice, restorableDevice uintptr)

var (
	mutex          sync.Mutex
	registedHandle []NotificationHander
	loopCancle     context.CancelFunc
	//mainLoopThreadMac uintptr
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
	/*
		if runtime.GOOS != "windows" {
			runtime.LockOSThread()
			mainLoopThreadMac = uintptr(C.CFRunLoopGetCurrent())
			C.CFRunLoopRun()
		}
	*/

	if loopCancle != nil {
		return
	}
	var contextForWin context.Context
	contextForWin, loopCancle = context.WithCancel(context.Background())
	select {
	case <-contextForWin.Done():
	}

}

//IsRunLoop ...
func IsRunLoop() bool {
	return loopCancle != nil
}

//RunLoopStop ...
func RunLoopStop() {
	/*
		if runtime.GOOS != "windows" {
			if 0 != mainLoopThreadMac {
				C.CFRunLoopStop(C.CFRunLoopRef(mainLoopThreadMac))
			}
		}
	*/
	if loopCancle != nil {
		loopCancle()
		loopCancle = nil
	}

}

//AMRestorableDeviceGetECID ...
func AMRestorableDeviceGetECID(restoredev uintptr) int64 {
	return int64(C.AMRestorableDeviceGetECID(unsafe.Pointer(restoredev)))
}

//AMDeviceIsPaired ...
func AMDeviceIsPaired(modeldev uintptr) int32 {
	return int32(C.AMDeviceIsPaired(unsafe.Pointer(modeldev)))
}

//AMDevicePair ...
func AMDevicePair(modeldev uintptr) int32 {
	return int32(C.AMDevicePair(unsafe.Pointer(modeldev)))
}

//AMDeviceValidatePairing ...
func AMDeviceValidatePairing(modeldev uintptr) int32 {
	return int32(C.AMDeviceValidatePairing(unsafe.Pointer(modeldev)))
}

//AMDeviceStartSession ...
func AMDeviceStartSession(modeldev uintptr) int32 {
	return int32(C.AMDeviceStartSession(unsafe.Pointer(modeldev)))
}

//AMDeviceStartService ...
func AMDeviceStartService(modeldev uintptr, name string) (connect uintptr, res uint32) {
	var nUnknown int
	res = uint32(C.AMDeviceStartService(
		unsafe.Pointer(modeldev),
		unsafe.Pointer(MakeCFString(name)),
		C.PPV(unsafe.Pointer(&connect)),
		unsafe.Pointer(&nUnknown)))
	return
}

//AMDServiceConnectionInvalidate ...
func AMDServiceConnectionInvalidate(connect uintptr) int {
	return int(C.AMDServiceConnectionInvalidate(unsafe.Pointer(connect)))
}

//AMDObserveNotification ...
func AMDObserveNotification(connect uintptr, name string) int32 {
	return int32(C.AMDObserveNotification(
		unsafe.Pointer(connect),
		unsafe.Pointer(MakeCFString(name))))
}

//AMDeviceStopSession ...
func AMDeviceStopSession(modedev uintptr) int32 {
	return int32(C.AMDeviceStopSession(
		unsafe.Pointer(modedev)))
}

//AMDeviceConnect ...
func AMDeviceConnect(modedev uintptr) int32 {
	return int32(C.AMDeviceConnect(
		unsafe.Pointer(modedev)))
}

//AMDeviceDisconnect ...
func AMDeviceDisconnect(modedev uintptr) int32 {
	return int32(C.AMDeviceDisconnect(
		unsafe.Pointer(modedev)))
}

//AMRestorableDeviceGetChipID ...
func AMRestorableDeviceGetChipID(restore uintptr) int {
	return int(C.AMRestorableDeviceGetChipID(
		unsafe.Pointer(restore)))
}

//AMRestorableDeviceGetBoardID ...
func AMRestorableDeviceGetBoardID(restore uintptr) int {
	return int(C.AMRestorableDeviceGetBoardID(
		unsafe.Pointer(restore)))
}

//AMRestoreModeDeviceGetTypeID ...
func AMRestoreModeDeviceGetTypeID(restore uintptr) uint {
	return uint(C.AMRestoreModeDeviceGetTypeID(
		unsafe.Pointer(restore)))
}

//AMRestorableDeviceGetProductID ...
func AMRestorableDeviceGetProductID(restore uintptr) uint {
	return uint(C.AMRestorableDeviceGetProductID(
		unsafe.Pointer(&restore)))
}

//AMRestorableDeviceGetProductType ...
func AMRestorableDeviceGetProductType(restore uintptr) uint {
	return uint(C.AMRestorableDeviceGetProductType(
		unsafe.Pointer(restore)))
}

//AMRestorableDeviceGetLocationID ...
func AMRestorableDeviceGetLocationID(restore uintptr) uint {
	return uint(C.AMRestorableDeviceGetLocationID(
		unsafe.Pointer(restore)))
}

//AMRestorableDeviceCopySerialNumber ...
func AMRestorableDeviceCopySerialNumber(restore uintptr) string {
	cf := C.AMRestorableDeviceCopySerialNumber(
		unsafe.Pointer(restore))
	if 0 != uintptr(cf) {
		defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(cf)))
	}
	return CFStringToString(uintptr(cf))
}

//AMRecoveryModeDeviceGetProductionMode ...
func AMRecoveryModeDeviceGetProductionMode(modeDevice uintptr) uint8 {
	return uint8(C.AMRecoveryModeDeviceGetProductionMode(
		unsafe.Pointer(modeDevice)))
}

//AMRecoveryModeDeviceGetTypeID ...
func AMRecoveryModeDeviceGetTypeID(modeDevice uintptr) uint {
	return uint(C.AMRecoveryModeDeviceGetTypeID(
		unsafe.Pointer(modeDevice)))
}

//AMDFUModeDeviceGetProductionMode ...
func AMDFUModeDeviceGetProductionMode(modeDevice uintptr) uint8 {
	return uint8(C.AMDFUModeDeviceGetProductionMode(
		unsafe.Pointer(modeDevice)))
}

//AMDFUModeDeviceGetTypeID ...
func AMDFUModeDeviceGetTypeID(modeDevice uintptr) uint {
	return uint(C.AMDFUModeDeviceGetTypeID(
		unsafe.Pointer(modeDevice)))
}

//AMDeviceCopyDeviceIdentifier ...
func AMDeviceCopyDeviceIdentifier(modeDevice uintptr) string {
	cf := C.AMDeviceCopyDeviceIdentifier(
		unsafe.Pointer(modeDevice))
	if 0 != uintptr(cf) {
		defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(cf)))
	}
	return CFStringToString(uintptr(cf))
}

//AMDeviceCopyValue ...
func AMDeviceCopyValue(modeDevice, v1, v2 uintptr) []byte {
	value := C.AMDeviceCopyValue(
		unsafe.Pointer(modeDevice),
		unsafe.Pointer(v1),
		unsafe.Pointer(v2))
	if 0 != uintptr(value) {
		defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(value)))
	}

	return CFToPlist(uintptr(unsafe.Pointer(value)))
}
