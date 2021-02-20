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

//AthServiceImpl ...
type AthServiceImpl struct {
	dev        mtunes.IOSDevice
	athConnect uintptr
	proxy      AthProxy
	status     string
}

//New ...
func New(dev mtunes.IOSDevice, proxy AthProxy) (ath AthService, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device is not trusted")
		return
	}

	ath = &AthServiceImpl{
		dev:    dev,
		proxy:  proxy,
		status: "status_allowed",
	}
	return
}

//Dial ...
func (ath *AthServiceImpl) Dial() (err error) {
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

//Exec ...
func (ath *AthServiceImpl) Exec(ctx context.Context) (err error) {
	msgChan := make(chan uintptr, 3)
	recvLoopCtx, cancelRecvFun := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ath.receiveLoop(recvLoopCtx, msgChan)
		wg.Done()
	}()

	outDispatch := make(chan struct{}, 5)

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
				err = ath.eventDispatch(ctx, outDispatch, msg)
				if err != nil {
					outDispatch <- struct{}{}
				}
			}()
		case <-ctx.Done():
			err = fmt.Errorf("Cancle")
			break dispatchFor
		case <-outDispatch:
			break dispatchFor
		}
	}
	cancelRecvFun()
	wg.Wait()
	return nil
}

func (ath *AthServiceImpl) receiveLoop(ctx context.Context, msgChan chan<- uintptr) {
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
			}
		}
		msgChan <- msg
	}
}

func (ath *AthServiceImpl) eventDispatch(ctx context.Context, recvFinish chan<- struct{}, msg uintptr) (err error) {
	//ignore message: Progress AssetMetrics InstalledAssets
	defer iapi.CFRelease(msg)

	msgName := iapi.ATCFMessageGetName(msg)
	if strings.EqualFold(msgName, "SyncFailed") {
		return fmt.Errorf("Recive syncFailed")
	}
	if strings.EqualFold(msgName, "ConnectionInvalid") {
		return fmt.Errorf("Recive ConnectionInvalid")
	}
	if strings.EqualFold(msgName, "Ping") {
		iapi.ATHostConnectionSendPing(ath.athConnect)
		return
	}

	if strings.EqualFold(msgName, "SyncFinished") {
		recvFinish <- struct{}{}
		return
	}

	if strings.EqualFold(msgName, "SyncAllowed") {
		if "status_allowed" != ath.status {
			return
		}
		ath.status = "status_ready"
		err = ath.allowEvent()
	} else if strings.EqualFold(msgName, "ReadyForSync") {
		if "status_ready" != ath.status {
			return
		}
		ath.status = "status_manifest"
		err = ath.readyEvent(ctx, msg)
	} else if strings.EqualFold(msgName, "AssetManifest") {
		if "status_manifest" != ath.status {
			return
		}
		err = ath.manifestEvent(ctx, msg)
	}
	return
}

func (ath *AthServiceImpl) allowEvent() (err error) {
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

func (ath *AthServiceImpl) readyEvent(ctx context.Context, msg uintptr) (err error) {
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

func (ath *AthServiceImpl) manifestEvent(ctx context.Context, msg uintptr) (err error) {
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

		r, size, err := ath.proxy.OpenEntityReader(assetID)
		if err != nil {
			continue
		}
		w, notify, err := ath.proxy.OpenEntityWriter(assetID)
		if err != nil {
			r.Close()
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
		} else {
			iapi.ATHostConnectionSendAssetCompleted(ath.athConnect, assetID, keybag, notify)
		}
	}
	return
}
