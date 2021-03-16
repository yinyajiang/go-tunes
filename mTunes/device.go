package mtunes

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	iapi "github.com/yinyajiang/go-tunes/itunesapi"
	tools "github.com/yinyajiang/go-ytools/utils"
)

const (
	serviceUserDataSuffix = "_#@Service_Conn@#_"
)

type deviceImpl struct {
	mode           string
	modeDevice     uintptr
	originalDevice uintptr
	info           map[string]interface{}
	buildSession   bool
	connected      bool

	extractCtx context.Context

	userLock sync.RWMutex
	userData map[string]interface{}
}

func newDeviceImpl(extractCtx context.Context, mode string, modeDevice, originalDevice uintptr) *deviceImpl {
	if len(mode) == 0 || modeDevice == 0 || originalDevice == 0 {
		return nil
	}
	dev := &deviceImpl{
		mode:           mode,
		modeDevice:     modeDevice,
		originalDevice: originalDevice,
		extractCtx:     extractCtx,
	}
	err := dev.loadBaseInfo()
	if err != nil {
		fmt.Println("Load base info fail,", err)
		return nil
	}

	//非拔出设备
	if extractCtx != nil && mode == "normal" {
		if nil != dev.Trust() { //try once
			dev.loadDetailInfo() //load nottrust info
		}
	}
	return dev
}

//SaveUserData ...
func (dev *deviceImpl) SaveUserData(key string, val interface{}) {
	dev.userLock.Lock()
	defer dev.userLock.Unlock()
	if dev.userData == nil {
		dev.userData = make(map[string]interface{}, 2)
	}
	dev.userData[key] = val
}

//GetUserData ...
func (dev *deviceImpl) GetUserData(key string) interface{} {
	dev.userLock.RLock()
	defer dev.userLock.RUnlock()
	if dev.userData == nil {
		return nil
	}
	return dev.userData[key]
}

//DeleteUserData ...
func (dev *deviceImpl) DeleteUserData(key string) {
	dev.userLock.Lock()
	defer dev.userLock.Unlock()
	if dev.userData == nil {
		return
	}
	delete(dev.userData, key)
}

//GetStartService ...
func (dev *deviceImpl) GetStartService(name string) (conn uintptr, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device not trust")
		return
	}

	conn, ok := dev.GetUserData(serviceUserDataKey(name)).(uintptr)
	if ok {
		return
	}

	conn, res := iapi.AMDeviceStartService(dev.ModeDevice(), name)
	if 0xe800001a == res || 0xE800001A == res { //kAMDPasswordProtectedError = 0xe800001a
		err = fmt.Errorf("Device password protected")
		return
	}

	if 0 == conn {
		err = fmt.Errorf("Start %s service fail", name)
		return
	}
	dev.SaveUserData(serviceUserDataKey(name), conn)
	return
}

func (dev *deviceImpl) Release() {
	dev.userLock.RLock()
	tmpData := dev.userData
	dev.userLock.RUnlock()
	for k := range tmpData {
		if !strings.HasSuffix(k, serviceUserDataSuffix) {
			continue
		}
		nm := k[0 : len(k)-len(serviceUserDataSuffix)]
		dev.StopService(nm)
	}
}

//IsServiceRuning ...
func (dev *deviceImpl) IsServiceRuning(name string) bool {
	_, ok := dev.GetUserData(serviceUserDataKey(name)).(uintptr)
	return ok
}

//StopService ...
func (dev *deviceImpl) StopService(name string) {
	conn, ok := dev.GetUserData(serviceUserDataKey(name)).(uintptr)
	if ok {
		if name != "com.apple.afc" {
			iapi.AMDServiceConnectionInvalidate(conn)
		}
		dev.DeleteUserData(serviceUserDataKey(name))
	}
}

//Mode ...
func (dev *deviceImpl) Mode() string {
	return dev.mode
}

//ID ...
func (dev *deviceImpl) ID() string {
	return dev.DeviceInfo()["id"].(string)
}

//ECID ...
func (dev *deviceImpl) ECID() string {
	return dev.DeviceInfo()["ecid"].(string)
}

//Model ...
func (dev *deviceImpl) Model() string {
	model, ok := dev.DeviceInfo()["model"].(string)
	if ok {
		return model
	}
	return ""
}

//Version ...
func (dev *deviceImpl) Version() string {
	ver, ok := dev.DeviceInfo()["ProductVersion"].(string)
	if ok {
		return ver
	}
	return ""

}

//ModelName ...
func (dev *deviceImpl) ModelName() string {
	modelName, ok := dev.DeviceInfo()["modelName"].(string)
	if ok {
		return modelName
	}
	return ""
}

//DeviceInfo ...
func (dev *deviceImpl) DeviceInfo() map[string]interface{} {
	dev.info["trust"] = dev.IsTrusted()
	return dev.info
}

//ModeDevice ...
func (dev *deviceImpl) ModeDevice() uintptr {
	return dev.modeDevice
}

//OriginalDevice ...
func (dev *deviceImpl) OriginalDevice() uintptr {
	return dev.originalDevice
}

//Trust ...
func (dev *deviceImpl) Trust() (err error) {
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

	connection, _ := dev.GetStartService("com.apple.mobile.notification_proxy")
	if connection != 0 {
		fmt.Println("AMDeviceStartService notification_proxy fail")
		iapi.AMDObserveNotification(connection, "com.apple.itunes-client.syncCancelRequest")
	}

	dev.buildSession = true

	dev.loadDetailInfo()
	return
}

//WaitTrust ...
func (dev *deviceImpl) WaitTrust(ctx context.Context) error {
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
func (dev *deviceImpl) AbordTrust() {
	if dev.buildSession {
		dev.buildSession = false
		iapi.AMDeviceStopSession(dev.ModeDevice())
		dev.disconnect()
	}
}

//IsTrusted ...
func (dev *deviceImpl) IsTrusted() bool {
	return dev.buildSession
}

//IsExtract ...
func (dev *deviceImpl) IsExtract() bool {
	if dev.extractCtx == nil {
		return true
	}
	select {
	case <-dev.extractCtx.Done():
		return true
	default:
	}
	return false
}

//ExtrackContext ...
func (dev *deviceImpl) ExtrackContext() context.Context {
	return dev.extractCtx
}

//ExtrackContext ...
func (dev *deviceImpl) WorkDir(join string) string {
	return tools.TempPath(tools.AbsJoinPath(dev.ID(), join))
}

func (dev *deviceImpl) loadBaseInfo() error {

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

func (dev *deviceImpl) loadDetailInfo() error {
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
func (dev *deviceImpl) connect() error {
	if "normal" != dev.Mode() {
		return fmt.Errorf("Connect not normal device")
	}
	if dev.connected {
		return nil
	}
	if 0 != iapi.AMDeviceConnect(dev.ModeDevice()) {
		return fmt.Errorf("Connect normal device fail")
	}
	dev.connected = true
	return nil
}

//Disconnect ...
func (dev *deviceImpl) disconnect() {
	if !dev.connected {
		return
	}
	if "normal" != dev.Mode() {
		return
	}
	dev.connected = false
	iapi.AMDeviceDisconnect(dev.ModeDevice())
}

func serviceUserDataKey(nm string) string {
	return nm + serviceUserDataSuffix
}
