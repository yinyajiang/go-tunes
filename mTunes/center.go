package mtunes

import (
	"context"
	"sync"
	"time"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
)

//NotifyEvent ...
type NotifyEvent struct {
	Device IOSDevice
	Event  string
}

var (
	devMutex          sync.Mutex
	devices           map[string]*IOSDeviceImpl     = make(map[string]*IOSDeviceImpl, 2)
	extractCancleFuns map[string]context.CancelFunc = make(map[string]context.CancelFunc, 2)
	subMutex          sync.Mutex
	subChans          map[chan *NotifyEvent]struct{} = make(map[chan *NotifyEvent]struct{}, 2)
	iapiSub           bool
)

//SubscriptionDeviceNotify ...
func SubscriptionDeviceNotify() (subChan <-chan *NotifyEvent) {
	return subscription(true)
}

//RunEventLoop ...
func RunEventLoop() {
	iapi.RunLoopRun()
}

//StopEventLoop ...
func StopEventLoop() {
	iapi.RunLoopStop()
}

//OnceEventLoopForWaitDevice 为指定设备启动一次事件循环
func OnceEventLoopForWaitDevice(ctx context.Context, id string) (dev IOSDevice) {
	subscription(false)
	dev = WaitForDevice(ctx, id)
	return
}

//UnsubscriptionDeviceNotify ..
func UnsubscriptionDeviceNotify(ch chan *NotifyEvent) {
	if nil == ch {
		return
	}
	subMutex.Lock()
	defer subMutex.Unlock()
	delete(subChans, ch)
	close(ch)
}

//DeviceCount ...
func DeviceCount() int {
	return len(devices)
}

//WaitForDevice ...
func WaitForDevice(ctx context.Context, id string) (dev IOSDevice) {
	for dev == nil {
		devMutex.Lock()
		devImpl, ok := devices[id]
		devMutex.Unlock()

		if ok {
			dev = devImpl
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
	return
}

//IsExistDevice ...
func IsExistDevice(id string) (b bool) {
	devMutex.Lock()
	defer devMutex.Unlock()
	_, b = devices[id]
	return
}

func subscription(retChan bool) (subChan <-chan *NotifyEvent) {
	if retChan {
		ch := make(chan *NotifyEvent, 5)
		subChan = ch

		subMutex.Lock()
		subChans[ch] = struct{}{}
		subMutex.Unlock()
	}

	if !iapiSub {
		iapi.AMRestorableDeviceRegisterForNotificationsForDevices(deviceEvent)
		iapiSub = true
	}

	return
}

func deviceEvent(even, mode string, modeDevice, restorableDevice uintptr) {
	var dev *IOSDeviceImpl
	devMutex.Lock()
	if even == "insert" {
		extractCtx, cancleFun := context.WithCancel(context.Background())
		dev = NewIOSDeviceImpl(extractCtx, mode, modeDevice, restorableDevice)
		devices[dev.ID()] = dev
		extractCancleFuns[dev.ID()] = cancleFun
	} else {
		dev = NewIOSDeviceImpl(nil, mode, modeDevice, restorableDevice)
		delDev, ok := devices[dev.ID()]
		if ok {
			extractCancleFuns[delDev.ID()]()
			delete(devices, delDev.ID())
			delete(extractCancleFuns, delDev.ID())
		}
	}
	devMutex.Unlock()

	subMutex.Lock()
	tempSubChans := subChans
	subMutex.Unlock()
	for c := range tempSubChans {
		select {
		case c <- &NotifyEvent{
			Event:  even,
			Device: dev,
		}:
		case <-time.After(time.Second * 2):
		}

	}
}
