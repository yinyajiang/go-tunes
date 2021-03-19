package athservice

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11

#include "cig.h"
*/
import "C"
import (
	"unsafe"

	iapi "github.com/yinyajiang/go-tunes/itunesapi"
)

//CigCalc 计算同步Plist的cig值
func CigCalc(grappa []byte, plist []byte) (cig []byte) {
	var cigbuff [21]byte
	ciglen := int32(21)
	if C.BOOL(1) != C.cigCalc(
		(*C.byte)(unsafe.Pointer(iapi.SpliceToPtr(grappa))),
		(*C.byte)(unsafe.Pointer(iapi.SpliceToPtr(plist))),
		C.int(len(plist)),
		(*C.byte)(unsafe.Pointer(iapi.SpliceToPtr(cigbuff[:]))),
		(*C.int)(unsafe.Pointer(&ciglen))) {
		return
	}
	return cigbuff[:ciglen]
}
