package mtunes

import (
	"context"
	"strconv"
	"sync"
	"time"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
)

//NotifyEvent ...
type NotifyEvent struct {
	Dev   Device
	Event string
}

var (
	devMutex          sync.Mutex
	devices           map[string]*deviceImpl        = make(map[string]*deviceImpl, 2)
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

//OnceForWaitDevice ...
func OnceForWaitDevice(ctx context.Context, id string) (dev Device) {
	subscription(false)
	dev = waitForDevice(ctx, id)
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

func waitForDevice(ctx context.Context, id string) (dev Device) {
	for dev == nil {
		devMutex.Lock()
		//转成10进制
		idInt, err := strconv.ParseInt(id, 0, 64)
		if err != nil {
			return nil
		}
		id = strconv.FormatInt(idInt, 10)
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
	var dev *deviceImpl
	devMutex.Lock()
	if even == "insert" {
		extractCtx, cancleFun := context.WithCancel(context.Background())
		dev = newDeviceImpl(extractCtx, mode, modeDevice, restorableDevice)
		devices[dev.ID()] = dev
		extractCancleFuns[dev.ID()] = cancleFun
	} else {
		dev = newDeviceImpl(nil, mode, modeDevice, restorableDevice)
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
			Event: even,
			Dev:   dev,
		}:
		case <-time.After(time.Second * 2):
		}

	}
}
