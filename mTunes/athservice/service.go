package athservice

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	tools "github.com/yinyajiang/go-ytools/utils"
)

type serviceImpl struct {
	dev        mtunes.Device
	athConnect uintptr
	proxy      AthProxy
	status     string
}

//New ...
func New(dev mtunes.Device, proxy AthProxy) (ath Service, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device is not trusted")
		return
	}

	ath = &serviceImpl{
		dev:    dev,
		proxy:  proxy,
		status: "status_allowed",
	}
	return
}

//Dial ...
func (ath *serviceImpl) Dial() (err error) {
	deviceIDface, ok := ath.dev.DeviceInfo()["uuid"]
	if !ok {
		return fmt.Errorf("deviceID is empty")
	}
	deviceID := deviceIDface.(string)
	ath.athConnect = iapi.ATHostConnectionCreateWithLibrary(getiTunesVersion(), deviceID)
	if 0 == ath.athConnect {
		return fmt.Errorf("Create connect fail")
	}
	iapi.ATHostConnectionSendPowerAssertion(ath.athConnect, true)

	req, err := requestProto()
	if err != nil {
		return
	}
	iapi.ATHostConnectionSendHostInfo(ath.athConnect, req)
	return
}

//Serve ...
func (ath *serviceImpl) Serve(ctx context.Context) (err error) {
	defer func() {
		if ath.athConnect != 0 {
			iapi.ATHostConnectionInvalidate(ath.athConnect)
			iapi.ATHostConnectionRelease(ath.athConnect)
			ath.athConnect = 0
		}
	}()

	msgChan := make(chan uintptr, 3)

	recvLoopCtx, cancelRecvFun := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ath.receiveLoop(recvLoopCtx, msgChan)
		wg.Done()
	}()

dispatchFor:
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				break dispatchFor
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer iapi.CFRelease(msg)
				evErr := ath.eventDispatch(ctx, cancelRecvFun, msg)
				if err == nil {
					err = evErr
				}
			}()
		case <-ctx.Done():
			err = fmt.Errorf("Cancle")
			break dispatchFor
		}
	}
	cancelRecvFun()
	wg.Wait()
	return
}

func (ath *serviceImpl) receiveLoop(ctx context.Context, msgChan chan<- uintptr) {
	lastSleep := time.Now().Unix()
	for {
		msg := iapi.ATHostConnectionReadMessage(ath.athConnect)
		if 0 == msg {
			break
		}

		select {
		case <-ctx.Done():
			close(msgChan)
			return
		default:
		}

		//减缓消息
		msgName := iapi.ATCFMessageGetName(msg)
		if strings.EqualFold(msgName, "Progress") {
			if time.Now().Unix()-lastSleep > 20 {
				time.Sleep(time.Millisecond * 1000)
				lastSleep = time.Now().Unix()
			}
		}
		msgChan <- msg
	}
}

func (ath *serviceImpl) eventDispatch(ctx context.Context, cancleFun context.CancelFunc, msg uintptr) (err error) {
	//ignore message: Progress AssetMetrics InstalledAssets

	defer func() {
		if err != nil {
			cancleFun()
		}
	}()

	msgName := iapi.ATCFMessageGetName(msg)

	switch msgName {
	case "SyncFailed":
		err = fmt.Errorf("Recive syncFailed")
	case "ConnectionInvalid":
		err = fmt.Errorf("Recive ConnectionInvalid")
	case "SyncFinished":
		cancleFun()
	case "Ping":
		iapi.ATHostConnectionSendPing(ath.athConnect)
	case "SyncAllowed":
		if "status_allowed" != ath.status {
			return
		}
		ath.status = "status_ready"
		err = ath.allowEvent()
	case "ReadyForSync":
		if "status_ready" != ath.status {
			return
		}
		ath.status = "status_manifest"
		err = ath.readyEvent(ctx, msg)
	case "AssetManifest":
		if "status_manifest" != ath.status {
			return
		}
		err = ath.manifestEvent(ctx, msg)
	}

	return
}

func (ath *serviceImpl) allowEvent() (err error) {
	keybag, anchors := ath.proxy.GetKeybag()
	if len(keybag) == 0 || len(anchors) == 0 {
		return fmt.Errorf("GetKeybag is empty")
	}
	pl, err := responseAllowedProto(keybag, anchors)
	sessionNum := iapi.ATHostConnectionGetCurrentSessionNumber(ath.athConnect)
	rspMsg := iapi.ATCFMessageCreate(sessionNum, "RequestingSync", pl)

	iapi.ATHostConnectionSendMessage(ath.athConnect, rspMsg)
	return
}

func (ath *serviceImpl) readyEvent(ctx context.Context, msg uintptr) (err error) {
	plmsg := iapi.CFToPlist(msg)
	keybag, _ := ath.proxy.GetKeybag()

	syncNum, grappa, err := unmarshalReadyProto(plmsg, keybag)
	if err != nil {
		return
	}

	err = ath.proxy.SubmitReadyPlist(syncNum, grappa)
	if err != nil {
		return
	}
	iapi.ATHostConnectionSendPowerAssertion(ath.athConnect, true)

	p1, p2, err := responseReadyProto(keybag, syncNum)
	if err != nil {
		return
	}

	iapi.ATHostConnectionSendMetadataSyncFinished(ath.athConnect, p1, p2)
	return
}

func (ath *serviceImpl) manifestEvent(ctx context.Context, msg uintptr) (err error) {
	plmsg := iapi.CFToPlist(msg)
	keybag, _ := ath.proxy.GetKeybag()

	assetArray, err := unmarshalManifestProto(plmsg, keybag)
	if err != nil {
		return
	}

	for _, assetID := range assetArray {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("Cancle when copy asset")
			return
		default:
		}
		if !ath.proxy.IsExistAsset(assetID) {
			iapi.ATHostConnectionSendFileError(ath.athConnect, assetID, keybag, 3)
			continue
		}
		r, size, openErr := ath.proxy.OpenAssetReader(assetID)
		if openErr != nil {
			iapi.ATHostConnectionSendFileError(ath.athConnect, assetID, keybag, 3)
			err = openErr
			ath.proxy.AssetFinish(assetID, false)
			continue
		}
		w, notify, openErr := ath.proxy.OpenAssetWriter(assetID)
		if openErr != nil {
			r.Close()
			iapi.ATHostConnectionSendFileError(ath.athConnect, assetID, keybag, 3)
			err = openErr
			ath.proxy.AssetFinish(assetID, false)
			continue
		}

		_, err = tools.CopyFun(ctx, size, w, r, func(total int64, prog float64) {
			iapi.ATHostConnectionSendFileProgress(ath.athConnect,
				assetID,
				keybag,
				0,
				int32(1067450368.0+5242880.0*prog),
				0)
		})
		w.Close()
		r.Close()
		if err != nil {
			iapi.ATHostConnectionSendFileError(ath.athConnect, assetID, keybag, 3)
			ath.proxy.AssetFinish(assetID, false)
		} else {
			iapi.ATHostConnectionSendAssetCompleted(ath.athConnect, assetID, keybag, notify)
			ath.proxy.AssetFinish(assetID, true)
		}
	}
	return
}
