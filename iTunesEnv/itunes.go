package itunes

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	tools "github.com/yinyajiang/go-ytools/utils"
)

//ITunes ...
type ITunes struct {
	Version          string
	MobileVersion    string
	InstallType      string
	corefundationDir string
	airtrafficDir    string
}

//New ...
func New() *ITunes {
	installType := ""
	if isInstallWin32() {
		installType = "win32"
	} else if isInstallUWP() {
		installType = "uwp"
	} else if runtime.GOOS == "darwin" {
		installType = "mac"
	}
	return &ITunes{
		Version:          version(),
		MobileVersion:    mobileDeviceVersion(),
		InstallType:      installType,
		corefundationDir: corefundationDir(),
		airtrafficDir:    airtrafficDir(),
	}
}

//IsDamaged ...
func (i *ITunes) IsDamaged() bool {
	if i.InstallType == "win32" {
		return !isServiceRuning()
	} else if i.InstallType == "uwp" {
		return len(i.MobileVersion) == 0 || !isServiceRuning()
	} else if i.InstallType == "mac" {
		return false
	}
	return true
}

//IsLowerVersion ...
func (i *ITunes) IsLowerVersion(modelName, model, iosver string) bool {
	return isLowerVersion(modelName, model, iosver, i.Version, i.MobileVersion)
}

//Repair ...
func (i *ITunes) Repair() bool {
	if i.InstallType == "win32" {
		startWin32Service()
	} else if i.InstallType == "uwp" {
		copyUWPEvnToDesk()
		startUWPService()
	} else if i.InstallType == "mac" {
		return true
	}
	return !i.IsDamaged()
}

//Auto ...
func (i *ITunes) Auto() bool {
	if i.InstallType == "mac" {
		return true
	}
	if len(i.InstallType) > 0 {
		if i.IsDamaged() {
			return i.Repair()
		}
		return true
	}

	return nil == i.Install(nil)
}

//Install ...
func (i *ITunes) Install(progFun func(Phase string, Prog float64)) (err error) {
	if i.InstallType == "mac" {
		return nil
	}
	jurl := tools.FrameConfigValue("iTunesUrl", getSysBit())
	if jurl == nil {
		return
	}
	url := jurl.MustString("")
	if len(url) == 0 {
		return
	}

	path := tools.TempPath("itunes/" + tools.PathName(url))
	err = tools.DownFileFun(url, path, func(total int64, prog float64) {
		if nil != progFun {
			progFun("download", prog)
		}
	})
	if err != nil {
		return
	}
	var wgroup sync.WaitGroup
	installFinish := false
	wgroup.Add(2)
	go func() {
		defer wgroup.Done()
		defer func() { installFinish = true }()
		depath := tools.AbsJoinPath(tools.AbsParent(path), "unzip")
		tools.RemovePath(depath)
		err = tools.DeCompressFun(path, depath, nil)
		if err != nil {
			return
		}
		pkgs := tools.FilterDeepFile(depath, []string{"*.msi"})
		if len(pkgs) == 0 {
			err = fmt.Errorf("Not Found *.msi file")
			return
		}
		for _, pkg := range pkgs {
			installMsi(pkg)
		}
	}()
	go func() {
		defer wgroup.Done()
		prog := float64(0)
		for !installFinish {
			prog += 3
			if prog >= 95 {
				prog = 95
			}
			select {
			case <-time.After(time.Second):
				{
					if nil != progFun {
						progFun("install", prog)
					}
				}

			}
		}
	}()
	wgroup.Wait()

	if err == nil {
		if nil != progFun {
			progFun("install", 100.0)
		}
	}

	t := New()
	*i = *t
	if !i.Repair() {
		err = fmt.Errorf("Install Fail")
	} else {
		fmt.Print("Install success")
	}
	return
}
