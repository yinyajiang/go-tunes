package itunes

import (
	"encoding/binary"
	"strconv"
	"strings"
	"unsafe"

	"github.com/yinyajiang/go-w32"
	"github.com/yinyajiang/go-w32/wutil"
)

func queryReg(path, key string, replacex86 bool) (reg string) {
	t, val := w32.RegGetRawAll(w32.HKEY_LOCAL_MACHINE, path, key)
	if t != 0 && len(val) > 0 {
		reg = w32.UTF16ByteToString(val)
	}

	if 64 == w32.GetSysBit() && replacex86 {
		if -1 == strings.Index(reg, "Program Files (x86)") {
			strings.ReplaceAll(reg, "Program Files", "Program Files (x86)")
		}
	}
	return
}

func setRegStartup(key string) {
	path := "Local Settings\\Software\\Microsoft\\Windows\\CurrentVersion\\AppModel\\SystemAppData"

	hMainKey := w32.RegOpenKeyEx(w32.HKEY_CLASSES_ROOT, path, w32.KEY_READ)
	if hMainKey != 0 {
		defer w32.RegCloseKey(hMainKey)
	}

	fullPath := ""
	for dwIndex := 0; ; dwIndex++ {
		name := w32.RegEnumKeyEx(hMainKey, uint32(dwIndex))
		if len(name) == 0 {
			break
		}
		if -1 != strings.Index(name, "AppleInc.iTunes") {
			fullPath = path + "\\" + name + "\\" + key
			break
		}
	}
	if len(fullPath) == 0 {
		return
	}

	hSubKey := w32.RegOpenKeyEx(w32.HKEY_CLASSES_ROOT, fullPath, w32.KEY_READ|w32.KEY_WRITE)
	if hSubKey != 0 {
		defer w32.RegCloseKey(hSubKey)
	}
	buf := make([]byte, unsafe.Sizeof(uint64(0)), unsafe.Sizeof(uint64(0)))
	binary.PutUvarint(buf, 2)
	w32.RegSetRaw(hSubKey, "State", w32.REG_DWORD, buf[:4])
}

func installMsi(pkg string) {
	_,h := wutil.StartAdminProcess("msiexec.exe", []string{"/i", `"` + pkg + `"`, "/qn", "/norestart"})
	w32.WaitForSingleObject(h, w32.INFINITE)
}

func getSysBit() string {
	return strconv.Itoa(w32.GetSysBit())
}
