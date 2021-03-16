package iapi

/*
#cgo windows CPPFLAGS: -DWIN32=1  -D_WIN32=1 -DUNICODE=1 -I.
#cgo CXXFLAGS: -std=c++11
#include "simpleApi.h"
*/
import "C"
import (
	"unsafe"

	tools "github.com/yinyajiang/go-ytools/utils"
)

//SetWinDllDir only for windows
func SetWinDllDir(path string) {
	C.AddLoadDir(C.PWCHAR(unsafe.Pointer(tools.StringToWCharPtr(path))))
}
