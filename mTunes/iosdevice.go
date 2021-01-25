package mtunes

import (
	"context"
	"fmt"
	"strconv"
	"time"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
)

//IOSDeviceImpl ...
type IOSDeviceImpl struct {
	mode           string
	modeDevice     uintptr
	originalDevice uintptr
	info           map[string]interface{}
	buildSession   bool
}

//NewIOSDevice ...
func NewIOSDevice(mode string, modeDevice, originalDevice uintptr) IOSDevice {
	if len(mode) == 0 || modeDevice == 0 || originalDevice == 0 {
		return nil
	}
	dev := &IOSDeviceImpl{
		mode:           mode,
		modeDevice:     modeDevice,
		originalDevice: originalDevice,
	}
	err := dev.loadBaseInfo()
	if err != nil {
		fmt.Println("Load base info fail,", err)
		return nil
	}
	if mode == "normal" {
		if nil != dev.Trust() { //try once
			dev.loadDetailInfo() //load nottrust info
		}
	}
	return dev
}

//Mode ...
func (dev *IOSDeviceImpl) Mode() string {
	return dev.mode
}

//ID ...
func (dev *IOSDeviceImpl) ID() string {
	return dev.DeviceInfo()["id"].(string)
}

//ECID ...
func (dev *IOSDeviceImpl) ECID() string {
	return dev.DeviceInfo()["ecid"].(string)
}

//Model ...
func (dev *IOSDeviceImpl) Model() string {
	model, ok := dev.DeviceInfo()["model"].(string)
	if ok {
		return model
	}
	return ""
}

//Version ...
func (dev *IOSDeviceImpl) Version() string {
	ver, ok := dev.DeviceInfo()["ProductVersion"].(string)
	if ok {
		return ver
	}
	return ""

}

//ModelName ...
func (dev *IOSDeviceImpl) ModelName() string {
	modelName, ok := dev.DeviceInfo()["modelName"].(string)
	if ok {
		return modelName
	}
	return ""
}

//DeviceInfo ...
func (dev *IOSDeviceImpl) DeviceInfo() map[string]interface{} {
	dev.info["trust"] = dev.IsTrusted()
	return dev.info
}

//ModeDevice ...
func (dev *IOSDeviceImpl) ModeDevice() uintptr {
	return dev.modeDevice
}

//OriginalDevice ...
func (dev *IOSDeviceImpl) OriginalDevice() uintptr {
	return dev.originalDevice
}

//Trust ...
func (dev *IOSDeviceImpl) Trust() (err error) {
	if dev.buildSession {
		return nil
	}
	if err = dev.connect(); err != nil {
		return
	}
	dis := true
	defer func() {
		if dis {
			dev.disconnect()
		}
	}()

	if 0 == iapi.AMDeviceIsPaired(dev.ModeDevice()) &&
		0 != iapi.AMDevicePair(dev.ModeDevice()) {
		return fmt.Errorf("Pair fail")
	}

	if 0 != iapi.AMDeviceValidatePairing(dev.ModeDevice()) {
		return fmt.Errorf("AMDeviceValidatePairing fail")
	}

	if 0 != iapi.AMDeviceStartSession(dev.ModeDevice()) {
		return fmt.Errorf("AMDeviceStartSession fail")
	}

	dis = false

	connection := iapi.AMDeviceStartService(dev.ModeDevice(), "com.apple.mobile.notification_proxy")
	if connection == 0 {
		return fmt.Errorf("AMDeviceStartService notification_proxy fail")
	}
	iapi.AMDObserveNotification(connection, "com.apple.itunes-client.syncCancelRequest")

	dev.buildSession = true

	dev.loadDetailInfo()
	return
}

//WaitTrust ...
func (dev *IOSDeviceImpl) WaitTrust(ctx context.Context) error {
	if "normal" != dev.Mode() {
		return fmt.Errorf("Is not normal mode")
	}

	for !dev.IsTrusted() {
		if dev.IsExtract() {
			return fmt.Errorf("Device extract")
		}
		if nil == dev.Trust() {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("Cancle or Timeout")
		case <-time.After(time.Second):
		}
	}
	return nil
}

//AbordTrust ...
func (dev *IOSDeviceImpl) AbordTrust() {
	if dev.buildSession {
		dev.buildSession = false
		iapi.AMDeviceStopSession(dev.ModeDevice())
		dev.disconnect()
	}
}

//IsTrusted ...
func (dev *IOSDeviceImpl) IsTrusted() bool {
	return dev.buildSession
}

//IsExtract ...
func (dev *IOSDeviceImpl) IsExtract() bool {
	return !IsExistDevice(dev.ID())
}

func (dev *IOSDeviceImpl) loadBaseInfo() error {

	ecid := iapi.AMRestorableDeviceGetECID(dev.OriginalDevice())
	if ecid == 0 {
		return fmt.Errorf("Get ECID fail")
	}

	dev.info = make(map[string]interface{}, 10)

	dev.info["ecid"] = strconv.FormatInt(ecid, 10)
	dev.info["id"] = strconv.FormatInt(ecid, 10)
	chipID := iapi.AMRestorableDeviceGetChipID(dev.OriginalDevice())
	boardID := iapi.AMRestorableDeviceGetBoardID(dev.OriginalDevice())
	dev.info["ChipID"] = chipID
	dev.info["BoardID"] = boardID

	dev.info["mode"] = dev.Mode()
	model := GetDeviceModel(int64(chipID), int64(boardID))
	if len(model) > 0 {
		dev.info["model"] = model
		dev.info["modelName"] = GetDeviceName(model)
		dev.info["type"] = GetDeviceType(model)
		dev.info["level"] = GetDeviceAppearanceLevel(model)
	}

	if "normal" == dev.Mode() {
		dev.info["TypeID"] = iapi.AMRestoreModeDeviceGetTypeID(dev.OriginalDevice())
		return nil
	}

	//recovery dfu
	dev.info["ProductID"] = iapi.AMRestorableDeviceGetProductID(dev.OriginalDevice())
	dev.info["ProductType"] = iapi.AMRestorableDeviceGetProductType(dev.OriginalDevice())
	dev.info["LocationID"] = iapi.AMRestorableDeviceGetLocationID(dev.OriginalDevice())
	dev.info["uuid"] = iapi.AMRestorableDeviceCopySerialNumber(dev.OriginalDevice())
	if "recovery" == dev.Mode() {
		dev.info["ProductionMode"] = iapi.AMRecoveryModeDeviceGetProductionMode(dev.ModeDevice())
		dev.info["TypeID"] = iapi.AMRecoveryModeDeviceGetTypeID(dev.ModeDevice())
	} else if "duf" == dev.Mode() {
		dev.info["ProductionMode"] = iapi.AMDFUModeDeviceGetProductionMode(dev.ModeDevice())
		dev.info["TypeID"] = iapi.AMDFUModeDeviceGetTypeID(dev.ModeDevice())
	}
	return nil
}

func (dev *IOSDeviceImpl) loadDetailInfo() error {
	if "normal" != dev.Mode() {
		return fmt.Errorf("Mode is not normal")
	}

	uuid := iapi.AMDeviceCopyDeviceIdentifier(dev.ModeDevice())
	if len(uuid) > 0 {
		dev.info["uuid"] = uuid
	}

	plinfo := iapi.AMDeviceCopyValue(dev.ModeDevice(), 0, 0)
	if len(plinfo) == 0 {
		return fmt.Errorf("AMDeviceCopyValue is empty")
	}
	var info map[string]interface{}
	UnmashalPlist(plinfo, &info)

	for key, val := range info {
		if "UniqueDeviceID" == key && len(uuid) == 0 {
			uuid = strconv.FormatInt(val.(int64), 10)
			dev.info["uuid"] = uuid
		} else if "ProductType" == key {
			model := val.(string)
			dev.info["model"] = model
		}
	}
	for key, val := range info {
		dev.info[key] = val
	}
	return nil
}

//Connect connect normal device
func (dev *IOSDeviceImpl) connect() error {
	if "normal" != dev.Mode() {
		return fmt.Errorf("Connect not normal device")
	}
	if 0 != iapi.AMDeviceConnect(dev.ModeDevice()) {
		return fmt.Errorf("Connect normal device fail")
	}
	return nil
}

//Disconnect ...
func (dev *IOSDeviceImpl) disconnect() {
	if "normal" != dev.Mode() {
		return
	}
	iapi.AMDeviceDisconnect(dev.ModeDevice())
}
