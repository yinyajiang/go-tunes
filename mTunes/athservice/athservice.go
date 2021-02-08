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

//AirService ...
type AirService struct {
	dev        mtunes.IOSDevice
	athConnect uintptr
	proxy      AirProxy
	status     string
}

//New ...
func New(dev mtunes.IOSDevice, proxy AirProxy) (ath *AirService, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device is not trusted")
		return
	}

	ath = &AirService{
		dev:    dev,
		proxy:  proxy,
		status: "status_allowed",
	}
	return
}

//Dial ...
func (ath *AirService) Dial() (err error) {
	deviceIDface, ok := ath.dev.DeviceInfo()["DeviceID"]
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
func (ath *AirService) Exec(ctx context.Context) (err error) {
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

func (ath *AirService) receiveLoop(ctx context.Context, msgChan chan<- uintptr) {
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

func (ath *AirService) eventDispatch(ctx context.Context, recvFinish chan<- struct{}, msg uintptr) (err error) {
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

func (ath *AirService) allowEvent() (err error) {
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

func (ath *AirService) readyEvent(ctx context.Context, msg uintptr) (err error) {
	stMsg := struct {
		Params struct {
			DataclassAnchors map[string]string `plist:"DataclassAnchors"`
			DeviceInfo       struct {
				Grappa []byte `plist:"Grappa"`
			} `plist:"DeviceInfo"`
		} `plist:"Params"`
	}{}

	plmsg := iapi.CFToPlist(msg)
	_, err = mtunes.UnmashalPlist(plmsg, &stMsg)
	if err != nil {
		return
	}
	keybag, _ := ath.proxy.GetKeybag()
	syncNum, ok := stMsg.Params.DataclassAnchors[keybag]
	if !ok {
		syncNum = stMsg.Params.DataclassAnchors["Media"]
	}
	if len(syncNum) == 0 {
		return fmt.Errorf("Not get sync num")
	}
	err = ath.proxy.SubmitReadyPlist(syncNum, stMsg.Params.DeviceInfo.Grappa)
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

func (ath *AirService) manifestEvent(ctx context.Context, msg uintptr) (err error) {

	type AssetItem struct {
		AssetID string `plist:"AssetID"`
	}
	type AssetArray []AssetItem
	stMsg := struct {
		Params struct {
			AssetManifest map[string]AssetArray `plist:"AssetManifest"`
		} `plist:"Params"`
	}{}

	plmsg := iapi.CFToPlist(msg)
	_, err = mtunes.UnmashalPlist(plmsg, &stMsg)
	if err != nil {
		return
	}
	if len(stMsg.Params.AssetManifest) == 0 {
		return fmt.Errorf("Mainifest event is empty")
	}
	keybag, _ := ath.proxy.GetKeybag()
	assetArray, ok := stMsg.Params.AssetManifest[keybag]
	if !ok {
		assetArray = stMsg.Params.AssetManifest["Media"]
	}
	if len(assetArray) == 0 {
		return fmt.Errorf("Mainifest event assetArray is empty")
	}

	for _, item := range assetArray {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("Cancle when copy asset")
			return
		default:
		}

		r, size, err := ath.proxy.OpenEntityReader(item.AssetID)
		if err != nil {
			continue
		}
		w, notify, err := ath.proxy.OpenEntityWriter(item.AssetID)
		if err != nil {
			r.Close()
			continue
		}
		_, err = tools.CopyFun(ctx, size, w, r, func(total int64, prog float64) {
			iapi.ATHostConnectionSendFileProgress(ath.athConnect,
				item.AssetID,
				keybag,
				0,
				int32(1067450368.0+5242880.0*prog),
				0)
		})
		w.Close()
		r.Close()
		if err != nil {
			iapi.ATHostConnectionSendFileError(ath.athConnect, item.AssetID, keybag, 3)
		} else {
			iapi.ATHostConnectionSendAssetCompleted(ath.athConnect, item.AssetID, keybag, notify)
		}
	}
	return
}
