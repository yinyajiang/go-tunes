package itunes

import (
	tools "github.com/yinyajiang/go-ytools/utils"
	"howett.net/plist"
)

func version() (ver string) {
	if !tools.IsExist("/Applications/iTunes.app/Contents/version.plist") {
		return
	}
	f, err := tools.OpenReadFile("/Applications/iTunes.app/Contents/version.plist")
	if err != nil {
		return
	}
	p := plist.NewDecoder(f)
	stver := struct {
		ShortVersionString string `plist:"CFBundleShortVersionString"`
		Version            string `plist:"CFBundleVersion"`
	}{}
	p.Decode(&stver)
	if len(stver.Version) > 0 {
		ver = stver.Version
	} else {
		ver = stver.ShortVersionString
	}
	return
}

func mobileDeviceVersion() string {
	return ""
}

func corefundationDir() (dir string) {
	return ""
}

func airtrafficDir() (dir string) {
	return ""
}

func isServiceRuning() bool {
	return true
}
func isLowerVersion(modelName, model, iosver, itunesver, mobilever string) bool {
	return false
}
func copyUWPEvnToDesk() {}

func startWin32Service() {}

func startUWPService() {}

func isInstallWin32() bool {
	return false
}
func isInstallUWP() bool {
	return false
}

func installMsi(pkg string) {}

func getSysBit() string {
	return "64"
}
