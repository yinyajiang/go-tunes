package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

#include "simpleApi.h"
*/
import "C"
import "unsafe"

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

//PlistToCF ...
func PlistToCF(data []byte) (cf uintptr) {
	if len(data) == 0 {
		return 0
	}
	cf = uintptr(C.MyPlistToCF(unsafe.Pointer(SpliceToPtr(data)), C.int(len(data))))
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

//MakeCFBool ...
func MakeCFBool(b bool) uintptr {
	if b {
		return uintptr(C.kCFBooleanTrue)
	}
	return uintptr(C.kCFBooleanFalse)
}

//CFBoolToBool ...
func CFBoolToBool(cfb uintptr) bool {
	return cfb == uintptr(C.kCFBooleanTrue)
}

//CFDictIsContainsKey ...
func CFDictIsContainsKey(cfdict uintptr, key string) bool {
	return 1 == int(C.MyCFDictIsContainsKey(unsafe.Pointer(cfdict), unsafe.Pointer(MakeCFString(key))))
}
