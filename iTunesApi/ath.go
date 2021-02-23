package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

#include "simpleApi.h"
*/
import "C"
import (
	"unsafe"
)

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
func ATHostConnectionSendPowerAssertion(athconn uintptr, b bool) int {
	return int(C.ATHostConnectionSendPowerAssertion(unsafe.Pointer(athconn),
		unsafe.Pointer(MakeCFBool(b))))
}

//ATHostConnectionRetain ...
func ATHostConnectionRetain(athconn uintptr) int {
	return int(C.ATHostConnectionRetain(unsafe.Pointer(athconn)))
}

//ATHostConnectionSendMetadataSyncFinished ...
func ATHostConnectionSendMetadataSyncFinished(athconn uintptr, plval1, plval2 []byte) int {
	return int(C.ATHostConnectionSendMetadataSyncFinished(unsafe.Pointer(athconn),
		unsafe.Pointer(PlistToCF(plval1)),
		unsafe.Pointer(PlistToCF(plval2))))
}

//ATHostConnectionSendFileError ...
func ATHostConnectionSendFileError(athconn uintptr, id, mtype string, val int32) {
	C.ATHostConnectionSendFileError(unsafe.Pointer(athconn),
		unsafe.Pointer(MakeCFString(id)),
		unsafe.Pointer(MakeCFString(mtype)),
		C.int(val))
}

//ATHostConnectionSendAssetCompleted ...
func ATHostConnectionSendAssetCompleted(athconn uintptr, id, mtype, path string) int {
	return int(C.ATHostConnectionSendAssetCompleted(unsafe.Pointer(athconn),
		unsafe.Pointer(MakeCFString(id)),
		unsafe.Pointer(MakeCFString(mtype)),
		unsafe.Pointer(MakeCFString(path))))
}

//ATCFMessageGetName ...
func ATCFMessageGetName(msg uintptr) string {
	cfname := uintptr(C.ATCFMessageGetName(unsafe.Pointer(msg)))
	if cfname == 0 {
		return ""
	}
	return CFStringToString(cfname)
}

//ATHostConnectionSendHostInfo ...
func ATHostConnectionSendHostInfo(athconn uintptr, plhost []byte) int {
	return int(C.ATHostConnectionSendHostInfo(unsafe.Pointer(athconn),
		unsafe.Pointer(PlistToCF(plhost))))
}

//ATHostConnectionSendMessage ...
func ATHostConnectionSendMessage(athconn, msg uintptr) int {
	return int(C.ATHostConnectionSendMessage(unsafe.Pointer(athconn),
		unsafe.Pointer(msg)))
}

//ATHostConnectionGetGrappaSessionId ...
func ATHostConnectionGetGrappaSessionId(athconn uintptr) int32 {
	return int32(C.ATHostConnectionGetGrappaSessionId(unsafe.Pointer(athconn)))
}

//ATHostConnectionReadMessage ...
func ATHostConnectionReadMessage(athconn uintptr) uintptr {
	return uintptr(C.ATHostConnectionReadMessage(unsafe.Pointer(athconn)))
}

//ATHostConnectionReadMessagePlist ...
func ATHostConnectionReadMessagePlist(athconn uintptr) (plmsg []byte) {
	msg := uintptr(C.ATHostConnectionReadMessage(unsafe.Pointer(athconn)))
	if msg == 0 {
		return
	}
	plmsg = CFToPlist(msg)
	C.CFRelease(C.CFTypeRef(msg))
	return
}
