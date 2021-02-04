package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

#include "simpleApi.h"
*/
import "C"
import "unsafe"

//ATCFMessageGetParam ...
func ATCFMessageGetParam(msg uintptr, name string) (para uintptr) {
	para = uintptr(C.ATCFMessageGetParam(unsafe.Pointer(msg),
		unsafe.Pointer(MakeCFString(name))))
	return
}

//ATHostConnectionGetCurrentSessionNumber ...
func ATHostConnectionGetCurrentSessionNumber(athconn uintptr) int32 {
	return int32(C.ATHostConnectionGetCurrentSessionNumber(unsafe.Pointer(athconn)))
}

//ATHostConnectionSendFileProgress ...
func ATHostConnectionSendFileProgress(athconn uintptr, v1, v2 string, v3 float32, v4 int32, v5 int32) {
	C.ATHostConnectionSendFileProgress(
		unsafe.Pointer(athconn),
		unsafe.Pointer(MakeCFString(v1)),
		unsafe.Pointer(MakeCFString(v2)),
		C.double(v3),
		C.int(v4),
		C.int(v5))
}

//ATCFMessageCreate ...
func ATCFMessageCreate(sessionNum int32, key string, plinfo []byte) (cfPara uintptr) {
	cfPara = uintptr(C.ATCFMessageCreate(C.int(sessionNum),
		unsafe.Pointer(MakeCFString(key)),
		unsafe.Pointer(PlistToCF(plinfo))))
	return
}

//ATHostConnectionCreateWithLibrary ...
func ATHostConnectionCreateWithLibrary(ituVer, deviceID string) uintptr {
	return uintptr(C.ATHostConnectionCreateWithLibrary(
		unsafe.Pointer(MakeCFString(ituVer)),
		unsafe.Pointer(MakeCFString(deviceID)),
		nil))
}

//ATHostConnectionSendPing ...
func ATHostConnectionSendPing(athconn uintptr) {
	C.ATHostConnectionSendPing(unsafe.Pointer(athconn))
}

//ATHostConnectionInvalidate ...
func ATHostConnectionInvalidate(athconn uintptr) {
	C.ATHostConnectionInvalidate(unsafe.Pointer(athconn))
}

//ATHostConnectionClose ...
func ATHostConnectionClose(athconn uintptr) {
	C.ATHostConnectionClose(unsafe.Pointer(athconn))
}

//ATHostConnectionRelease ...
func ATHostConnectionRelease(athconn uintptr) {
	C.ATHostConnectionRelease(unsafe.Pointer(athconn))
}

//ATHostConnectionSendPowerAssertion ...
func ATHostConnectionSendPowerAssertion(athconn uintptr, b bool) {
	if b {
		C.ATHostConnectionSendPowerAssertion(unsafe.Pointer(athconn),
			unsafe.Pointer(C.kCFBooleanTrue))
	}
}