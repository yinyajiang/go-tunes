package itunes

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yinyajiang/go-w32"
	"github.com/yinyajiang/go-w32/wutil"
	tools "github.com/yinyajiang/go-ytools/utils"
)

const (
	elem32 = `
[
	{
		"recursion": true,
		"srcs": [
			"{%UWPItunes%}/AMDS32/*",
			"{%UWPItunes%}/CFNetwork.resources",
			"{%UWPItunes%}/Foundation.resources",
			"{%UWPItunes%}/CoreFoundation.resources",
			"{%UWPItunes%}/CoreFoundation.dll",
			"{%UWPItunes%}/distnoted.exe"
		],
		"dstDir": "{%DeskItunes%}/Mobile Device Support"
	},
	{
		"recursion": true,
		"srcs": [
			"{%UWPItunes%}/CFNetwork.resources",
			"{%UWPItunes%}/Foundation.resources",
			"{%UWPItunes%}/CoreFoundation.resources"
		],
		"dstDir": "{%DeskItunes%}/Apple Application Support"
	},
	{
		"recursion": false,
		"srcs": [
			"{%UWPItunes%}/*",
			"{%UWPItunes%}/VFS/SystemX86/*"
		],
		"dstDir": "{%DeskItunes%}/Apple Application Support",
		"exts": [
			"*.dll",
			"*.exe"
		],
		"exclues": []
	}
]
`

	elem64 = `
[
	{
		"recursion": true,
		"srcs": [
			"{%UWPItunes%}/VFS/ProgramFilesCommonX86/Apple/Apple Application Support/*",
			"{%UWPItunes%}/VFS/SystemX86/*"
		],
		"dstDir": "{%DeskItunes%}/Apple Application Support"
	},
	{
		"recursion": true,
		"srcs": [
			"{%UWPItunes%}/AMDS32/*",
			"{%UWPItunes%}/VFS/ProgramFilesCommonX86/Apple/Apple Application Support/CFNetwork.resources",
			"{%UWPItunes%}/VFS/ProgramFilesCommonX86/Apple/Apple Application Support/Foundation.resources",
			"{%UWPItunes%}/VFS/ProgramFilesCommonX86/Apple/Apple Application Support/CoreFoundation.resources"
		],
		"dstDir": "{%DeskItunes%}/Mobile Device Support"
	}
]
`
	regTemplate = `
Windows Registry Editor Version 5.00
[HKEY_LOCAL_MACHINE\SOFTWARE{%WoW6432%}Apple Inc.]
[HKEY_LOCAL_MACHINE\SOFTWARE{%WoW6432%}Apple Inc.\Apple Application Support]
"InstallDir"="{%PFCF86%}\Apple\Apple Application Support"
"UserVisibleVersion"="7.0.2"
"Version"="7.0.2"
"AddByMobimover" = "1"

[HKEY_LOCAL_MACHINE\SOFTWARE{%WoW6432%}Apple Inc.\Apple Mobile Device Support]
[HKEY_LOCAL_MACHINE\SOFTWARE{%WoW6432%}Apple Inc.\Apple Mobile Device Support\Shared]
"ASMapiInterfaceDLL"="{%PFCF86%}\Apple\Mobile Device Support\AppleSyncMapiInterface.dll"
"MobileDeviceDLL"="{%PFCF86%}\Apple\Mobile Device Support\MobileDevice.dll"
"AirTrafficHostDLL"="{%PFCF86%}\Apple\Mobile Device Support\AirTrafficHost.dll"
"AddByMobimover" = "1"

`
)

var (
	str32Replace = map[string]string{
		"{%WoW6432%}":    `\`,
		"{%PFCF86%}":     `#43`,
		"{%DeskItunes%}": `#43/Apple`,
	}

	str64Replace = map[string]string{
		"{%WoW6432%}":    `\WOW6432Node\`,
		"{%PFCF86%}":     `#44`,
		"{%DeskItunes%}": `#44/Apple`,
	}
)

//GetUWPRegTemplate ..
func getUWPRegTemplate() string {
	return tools.ReplaceString(regTemplate, getReplaceMap())
}

//GetUWPCopyElem ...
func getUWPCopyElem() (elems []Elem) {
	var origi string
	if 32 == w32.GetSysBit() {
		origi = elem32
	} else {
		origi = elem64
	}
	origi = tools.ReplaceString(origi, getReplaceMap())
	err := json.Unmarshal([]byte(origi), &elems)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func getReplaceMap() (ret map[string]string) {
	ret = make(map[string]string, 0)
	b, desc := wutil.LoadUWPDesc("AppleInc.iTunes")
	if b {
		ret["{%UWPItunes%}"] = desc.WorkDir
	}

	var origi map[string]string
	if 32 == w32.GetSysBit() {
		origi = str32Replace
	} else {
		origi = str64Replace
	}
	reg := regexp.MustCompile(`\d+`)
	for k, v := range origi {
		if strings.HasPrefix(v, "#") {
			strSpecil := reg.FindString(v)
			specil, _ := strconv.Atoi(strSpecil)
			v = strings.ReplaceAll(v, "#"+strSpecil, w32.SHGetSpecialFolderPath(int32(specil)))
		}
		ret[k] = v
	}

	return
}
