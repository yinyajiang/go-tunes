package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

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
func AMDeviceStartService(modeldev uintptr, name string) (connect uintptr) {
	var nUnknown int
	C.AMDeviceStartService(
		unsafe.Pointer(modeldev),
		unsafe.Pointer(MakeCFString(name)),
		C.PPV(unsafe.Pointer(&connect)),
		unsafe.Pointer(&nUnknown))
	return
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

//MakeCFString ...
func MakeCFString(str string) uintptr {
	cstr := C.CString(str)
	if 0 == uintptr(unsafe.Pointer(cstr)) {
		return 0
	}
	defer C.free(unsafe.Pointer(cstr))
	return uintptr(unsafe.Pointer(C.CFStringCreateWithCString(C.CFAllocatorRef(unsafe.Pointer(uintptr(0))), cstr, C.kCFStringEncodingUTF8)))
}

//CFToPlist ...
func CFToPlist(cf uintptr) (data []byte) {
	if cf == 0 {
		return nil
	}
	var len int32
	plistbuff := C.MyCFToPlist(unsafe.Pointer(cf), (*C.int)(unsafe.Pointer(&len)))
	if plistbuff == nil || len == 0 {
		return nil
	}
	data = (*[1 << 30]byte)(unsafe.Pointer(plistbuff))[0:len:len]
	return
}

//SpliceToPtr ...
func SpliceToPtr(data []byte) uintptr {
	if len(data) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&data[0]))
}

//CFStringToString ...
func CFStringToString(cfstr uintptr) string {
	if 0 == cfstr {
		return ""
	}
	cfref := C.CFStringRef(unsafe.Pointer(cfstr))
	len := C.CFStringGetLength(cfref)
	if len <= 0 {
		return ""
	}
	buff := make([]byte, len, len+1)
	C.CFStringGetCString(cfref, (*C.char)(unsafe.Pointer(SpliceToPtr(buff))), len+1, C.kCFStringEncodingUTF8)
	return string(buff)
}
