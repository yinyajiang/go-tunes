package itunes

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	tools "github.com/yinyajiang/go-ytools/utils"

	"github.com/yinyajiang/go-w32/wutil"
)

//IsInstallWin32 ...
func isInstallWin32() bool {
	apps := wutil.LoadAppList()
	for name, val := range apps {
		if -1 != strings.Index(name, "Apple Mobile Device Support") {
			return true
		}
		if -1 != strings.Index(val, "Mobile Device Support") {
			return true
		}
	}
	return false
}

//IsInstallUWP ...
func isInstallUWP() bool {
	b, desc := wutil.LoadUWPDesc("AppleInc.iTunes")
	if !b {
		return false
	}
	return len(desc.WorkDir) > 0
}

//IsServiceRuning ...
func isServiceRuning() bool {
	return wutil.IsProcessRuning("AppleMobileDeviceProcess.exe", "AppleMobileDeviceService.exe", "MobileDeviceProcess.exe", "MobileDeviceService.exe")
}

//StartWin32Service ...
func startWin32Service() {
	cmd := exec.Command("sc", "config", `"Apple Mobile Device Service"`, "start=auto")
	if nil != cmd {
		cmd.Wait()
	}
	cmd = exec.Command("sc", "start", `"Apple Mobile Device Service"`)
	if nil != cmd {
		cmd.Wait()
	}
}

//StartUWPService ...
func startUWPService() {
	b, desc := wutil.LoadUWPDesc("AppleInc.iTunes")
	if !b || !tools.IsExist(desc.ExePath) {
		return
	}
	iTunesExe := strings.ReplaceAll(desc.ExePath, "\\", "/")

	err := wutil.DropPrivilegeStartProcess(iTunesExe)
	if err != nil {
		return
	}

	for i := 5; !isServiceRuning() && i > 0; i-- {
		select {
		case <-time.After(time.Second * 1):
		}
	}
}

//CorefundationDir ...
func corefundationDir() (dir string) {
	dir = queryReg("SOFTWARE\\Apple Inc.\\Apple Mobile Device Support", "InstallDir", true)
	if len(dir) > 0 && tools.IsExist(tools.ThePath(dir, "CoreFoundation.dll")) {
		return
	}

	dll := queryReg("SOFTWARE\\Apple Inc.\\Apple Mobile Device Support\\Shared", "ASMapiInterfaceDLL", false)
	dir = tools.AbsParent(dll)
	if len(dir) > 0 && tools.IsExist(tools.ThePath(dir, "CoreFoundation.dll")) {
		return
	}

	dir = queryReg("SOFTWARE\\Apple Inc.\\Apple Application Support", "InstallDir", true)
	if len(dir) > 0 && tools.IsExist(tools.ThePath(dir, "CoreFoundation.dll")) {
		return
	}
	dir = ""
	return
}

//AirtrafficDir ...
func airtrafficDir() (dir string) {
	dll := queryReg("SOFTWARE\\Apple Inc.\\Apple Mobile Device Support\\Shared", "AirTrafficHostDLL", true)
	if !tools.IsExist(dll) {
		dll = queryReg("SOFTWARE\\Apple Inc.\\Apple Mobile Device Support\\Shared", "ASMapiInterfaceDLL", false)
		dll = tools.ThePath(tools.AbsParent(dll), "AirTrafficHost.dll")
		if tools.IsExist(dll) {
			dir = tools.AbsParent(dll)
		}
	} else {
		dir = tools.AbsParent(dll)
	}
	return
}

//Version ...
func version() (ver string) {
	b, desc := wutil.LoadUWPDesc("AppleInc.iTunes")
	if !b {
		return ""
	}
	if len(desc.WorkDir) != 0 {
		ver = wutil.GetFileVersion(tools.ThePath(desc.WorkDir, "iTunes.exe"))
		if len(ver) != 0 {
			return
		}
	}
	ver = queryReg("SOFTWARE\\Apple Computer, Inc.\\iTunes", "Version", true)
	return
}

//MobileDeviceVersion ...
func mobileDeviceVersion() (ver string) {
	ver = queryReg("SOFTWARE\\Apple Inc.\\Apple Mobile Device Support", "Version", true)
	return
}

//IsLowerVersion ...
func isLowerVersion(modelName, model, iosver, itunesver, mobilever string) bool {
	cmpdev := "iPhone11,1"
	modelName = strings.ToLower(modelName)
	if 0 == len(itunesver) {
		itunesver = version()
	}
	if 0 == len(mobilever) {
		mobilever = mobileDeviceVersion()
	}
	if -1 != strings.Index(modelName, "ipad") {
		if -1 != strings.Index(modelName, "mini") {
			cmpdev = "iPad5,2"
		} else if -1 != strings.Index(modelName, "pro") {
			cmpdev = "iPad7,1"
		} else if -1 != strings.Index(modelName, "air") {
			cmpdev = "iPad5,4"
		} else {
			cmpdev = "iPad7,6"
		}
	}

	highIosVer := tools.CmpVersion(iosver, "13.6") >= 0
	newDevice := tools.CmpVersion(model, cmpdev) > 0
	if !highIosVer && !newDevice {
		return false
	}

	cmpItunes := ""
	cmpMobile := ""
	if runtime.GOOS == "windows" {
		if highIosVer || newDevice {
			cmpItunes = "12.10.10.2"
			cmpMobile = "14.1.0.35"
		} else {
			cmpItunes = "12.9.0"
			cmpMobile = "12.0.0.1039"
		}
	} else {
		cmpItunes = "12.8.0"
	}

	//优先使用MobileDeviceVersion
	insVersion := itunesver
	cmpVer := cmpItunes
	if 0 != len(mobilever) {
		insVersion = mobilever
		cmpVer = cmpMobile
	}
	if len(insVersion) == 0 {
		return false
	}
	return tools.CmpVersion(insVersion, cmpVer) < 0
}

//CopyUWPEvnToDesk ...
func copyUWPEvnToDesk() {
	//copy file
	elems := getUWPCopyElem()
	if len(elems) == 0 {
		return
	}
	copyElem(elems)
	//copy reg
	reg := getUWPRegTemplate()
	if len(reg) == 0 {
		return
	}
	path := tools.AbsJoinPath(os.TempDir(), strconv.Itoa(int(tools.RandNum()))+".reg")
	if nil != tools.WriteFileString(path, reg) {
		return
	}
	cmd := exec.Command("regedit.exe", "/s", `"`+path+`"`)
	if cmd == nil {
		return
	}
	cmd.Run()
	//set startup
	setRegStartup("AppleMobileDeviceProcess")
	setRegStartup("iTunesHelper")
}
