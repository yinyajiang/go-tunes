package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#cgo !windows LDFLAGS: -framework CoreFoundation -ObjC

#include "simpleApi.h"
*/
import "C"
import "unsafe"

//AFCConnectionOpen ...
func AFCConnectionOpen(afcconn uintptr) (conn uintptr) {
	C.AFCConnectionOpen(
		unsafe.Pointer(afcconn),
		C.int(0),
		C.PPV(unsafe.Pointer(&conn)))
	return
}

//AFCConnectionClose ...
func AFCConnectionClose(conn uintptr) int {
	return int(C.AFCConnectionClose(unsafe.Pointer(conn)))
}

//AFCDeviceInfoOpen ...
func AFCDeviceInfoOpen(conn uintptr) (filedict uintptr) {
	C.AFCDeviceInfoOpen(unsafe.Pointer(conn), C.PPV(unsafe.Pointer(&filedict)))
	return
}

//AFCFileInfoOpen ...
func AFCFileInfoOpen(conn uintptr, path string) (hand uintptr) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.AFCFileInfoOpen(unsafe.Pointer(conn),
		unsafe.Pointer(cpath),
		C.PPV(unsafe.Pointer(&hand)))
	return
}

//AFCKeyValueRead ...
func AFCKeyValueRead(infohand uintptr) (key string, value string) {
	var szkey, szvalue *C.char
	res := int(C.AFCKeyValueRead(unsafe.Pointer(infohand),
		C.PPV(unsafe.Pointer(&szkey)),
		C.PPV(unsafe.Pointer(&szvalue))))
	if 0 == res {
		key = C.GoString(szkey)
		value = C.GoString(szvalue)
	}
	return
}

//AFCKeyValueClose ...
func AFCKeyValueClose(infohand uintptr) int {
	return int(C.AFCKeyValueClose(unsafe.Pointer(infohand)))
}

//AFCDirectoryOpen ...
func AFCDirectoryOpen(conn uintptr, path string) (dirhand uintptr) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.AFCDirectoryOpen(unsafe.Pointer(conn),
		unsafe.Pointer(cpath),
		C.PPV(unsafe.Pointer(&dirhand)))
	return
}

//AFCDirectoryRead ...
func AFCDirectoryRead(conn uintptr, dirhand uintptr) (value string) {
	var szvalue *C.char
	res := int(C.AFCDirectoryRead(unsafe.Pointer(conn),
		unsafe.Pointer(dirhand),
		C.PPV(unsafe.Pointer(&szvalue))))
	if 0 == res {
		value = C.GoString(szvalue)
	}
	return
}

//AFCDirectoryClose ...
func AFCDirectoryClose(conn uintptr, dirhand uintptr) int {
	return int(C.AFCDirectoryClose(unsafe.Pointer(conn),
		unsafe.Pointer(dirhand)))
}

//AFCDirectoryCreate ...
func AFCDirectoryCreate(conn uintptr, path string) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	return int(C.AFCDirectoryCreate(unsafe.Pointer(conn),
		unsafe.Pointer(cpath)))
}

//AFCRemovePath ...
func AFCRemovePath(conn uintptr, path string) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	return int(C.AFCRemovePath(unsafe.Pointer(conn),
		unsafe.Pointer(cpath)))
}

//AFCRenamePath ...
func AFCRenamePath(conn uintptr, oldpath, newpath string) int {
	coldpath := C.CString(oldpath)
	defer C.free(unsafe.Pointer(coldpath))

	cnewpath := C.CString(newpath)
	defer C.free(unsafe.Pointer(cnewpath))

	return int(C.AFCRenamePath(unsafe.Pointer(conn),
		unsafe.Pointer(coldpath),
		unsafe.Pointer(cnewpath)))
}

//AFCFileRefOpen ...
func AFCFileRefOpen(conn uintptr, path string, mode int64) (fileHand uint64, res int32) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	res = int32(C.AFCFileRefOpen(unsafe.Pointer(conn),
		unsafe.Pointer(cpath),
		C.ulonglong(mode),
		(*C.ulonglong)(unsafe.Pointer(&fileHand))))
	return
}

//AFCFileRefRead ...
func AFCFileRefRead(conn uintptr, filehand uint64, buff []byte) (readed int, res int) {
	size := uint32(len(buff))
	res = int(C.AFCFileRefRead(unsafe.Pointer(conn),
		C.ulonglong(filehand),
		unsafe.Pointer(SpliceToPtr(buff)),
		(unsafe.Pointer(&size))))
	readed = int(size)
	return
}

//AFCFileRefWrite ...
func AFCFileRefWrite(conn uintptr, filehand uint64, data []byte) int {
	return int(C.AFCFileRefWrite(unsafe.Pointer(conn),
		C.ulonglong(filehand),
		unsafe.Pointer(SpliceToPtr(data)),
		C.int(len(data))))
}

//AFCFileRefClose ...
func AFCFileRefClose(conn uintptr, filehand uint64) int {
	return int(C.AFCFileRefClose(unsafe.Pointer(conn),
		C.ulonglong(filehand)))
}

//AFCFileRefSeek ...
func AFCFileRefSeek(conn uintptr, filehand uint64, seek uint64, mode int) int {
	return int(C.AFCFileRefSeek(unsafe.Pointer(conn),
		C.ulonglong(filehand),
		C.ulonglong(seek),
		C.ulong(mode)))
}

//AFCFileRefTell ...
func AFCFileRefTell(conn uintptr, filehand uint64) (seek uint64, res int) {
	res = int(C.AFCFileRefTell(unsafe.Pointer(conn),
		C.ulonglong(filehand),
		(*C.ulonglong)(unsafe.Pointer(&seek))))
	return
}
