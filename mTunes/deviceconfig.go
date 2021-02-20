package mtunes

import (
	"encoding/json"
	"fmt"
	"strings"

	tools "github.com/yinyajiang/go-ytools/utils"
)

//DeviceConfigInfo ...
type DeviceConfigInfo struct {
	Name        string `json:"name"`
	Identifier  string `json:identifier"`
	Boardconfig string `json:"boardconfig"`
	Platform    string `json:"platform"`
	Cpid        int64  `json:"cpid"`
	Bdid        int64  `json:"bdid"`
}

var deviceList []DeviceConfigInfo

//GetDeviceName ...
func GetDeviceName(model string) string {
	for _, info := range GetDeviceConfigList() {
		if info.Identifier == model {
			return info.Name
		}
	}
	return ""
}

//GetDeviceType ...
func GetDeviceType(model string) string {
	if strings.HasPrefix(model, "iPhone") {
		return "iPhone"
	} else if strings.HasPrefix(model, "iPad") {
		return "iPad"
	} else if strings.HasPrefix(model, "iPod") {
		return "iPod"
	}
	return "unknow"
}

//GetDeviceModel ...
func GetDeviceModel(chipID, boardID int64) string {
	for _, info := range GetDeviceConfigList() {
		if info.Cpid == chipID && info.Bdid == boardID {
			return info.Identifier
		}
	}
	return "unknow"
}

//GetDeviceAppearanceLevel ...
func GetDeviceAppearanceLevel(model string) int {
	switch GetDeviceType(model) {
	case "iPhone":
		if tools.CmpVersion(model, "iPhone12,8") == 0 { //iPhone se2
			return 0
		} else if tools.CmpVersion(model, "iPhone10,3") >= 0 { //iPhone x
			return 3
		} else if tools.CmpVersion(model, "iPhone10,1") >= 0 {
			return 2
		} else if tools.CmpVersion(model, "iPhone9,1") >= 0 &&
			tools.CmpVersion(model, "iPhone9,4") <= 0 {
			return 1
		}
	case "iPad":
		if tools.CmpVersion(model, "iPad8,1") >= 0 &&
			tools.CmpVersion(model, "iPad8,12") <= 0 {
			return 1
		}
	case "iPod Touch":
		if tools.CmpVersion(model, "iPod9,1") >= 0 {
			return 1
		}
	}
	return 0
}

//GetDeviceConfigList ...
func GetDeviceConfigList() []DeviceConfigInfo {
	if deviceList != nil {
		return deviceList
	}
	data, err := tools.ReadFileAll(tools.LocalPath("devices.json"))
	if err != nil {
		fmt.Printf("Read devices.json fail:%v", err)
		return nil
	}
	err = json.Unmarshal(data, &deviceList)
	if err != nil {
		fmt.Printf("Unmarshal devices.json fail:%v", err)
		return nil
	}
	return deviceList
}
