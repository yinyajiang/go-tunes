package ringtone

import (
	"fmt"
	"io"
	"strconv"

	"github.com/yinyajiang/go-tunes/service/athservice"
	"github.com/yinyajiang/go-tunes/service/fileservice"
	tools "github.com/yinyajiang/go-ytools/utils"
)

type syncProxy struct {
	inserts map[uint64]*TrackInfo
	dels    map[uint64]*TrackInfo
	fs      fileservice.Service
}

func (p *syncProxy) AssetFinish(entityID string, successed bool) {

}

//GetKeybag ...
func (p *syncProxy) GetKeybag() (keybag string, anchors []string) {
	return "Ringtone", []string{"0"}
}

//OpenAssetWriter ...
func (p *syncProxy) OpenAssetWriter(entityID string) (w io.WriteCloser, target string, err error) {
	pid, err := strconv.ParseUint(entityID, 0, 64)
	if err != nil {
		return
	}
	if p.inserts == nil {
		err = fmt.Errorf("inserts map is empty")
		return
	}
	track, ok := p.inserts[pid]
	if !ok {
		err = fmt.Errorf("inserts item not found")
		return
	}
	if !track.isFake {
		err = fmt.Errorf("inserts item is not fake")
		return
	}
	w, err = p.fs.OpenFile(track.Path, fileservice.AFC_FOPEN_WRONLY)
	target = track.Path
	return
}

//IsExistEntity ...
func (p *syncProxy) IsExistAsset(entityID string) bool {
	pid, err := strconv.ParseUint(entityID, 0, 64)
	if err != nil {
		return false
	}
	_, ok := p.inserts[pid]
	return ok
}

//OpenAssetReader ...
func (p *syncProxy) OpenAssetReader(entityID string) (r io.ReadCloser, size int64, err error) {
	pid, err := strconv.ParseUint(entityID, 0, 64)
	if err != nil {
		return
	}
	if p.inserts == nil {
		err = fmt.Errorf("inserts map is empty")
		return
	}
	track, ok := p.inserts[pid]
	if !ok {
		err = fmt.Errorf("inserts item not found")
		return
	}
	if !track.isFake {
		err = fmt.Errorf("inserts item is not fake")
		return
	}
	size = track.Size
	if track.dev == nil {
		r, err = tools.OpenReadFile(track.fakeSrcPath)
	} else {
		r, err = p.fs.OpenFile(track.fakeSrcPath, fileservice.AFC_FOPEN_RDONLY)
	}

	return
}

//SubmitReadyPlist ...
func (p *syncProxy) SubmitReadyPlist(strSyncNum string, grapa []byte) (err error) {
	insertArr := make([]*TrackInfo, 0, len(p.inserts))
	delArr := make([]*TrackInfo, 0, len(p.dels))
	for _, t := range p.inserts {
		insertArr = append(insertArr, t)
	}
	for _, t := range p.dels {
		delArr = append(delArr, t)
	}

	syncNum, _ := strconv.Atoi(strSyncNum)
	pl, err := mashalUpdateProto(syncNum, insertArr, delArr)
	if err != nil {
		return
	}
	path := fmt.Sprintf("/iTunes_Control/Ringtones/Sync/Sync_%08d.plist", uint(syncNum))
	if err != nil {
		return
	}
	err = p.writeSyncPlist(path, nil, pl)
	if err != nil {
		return
	}
	path = fmt.Sprintf("/iTunes_Control/Sync/Media/Sync_%08d.plist", uint(syncNum))
	if err != nil {
		return
	}
	err = p.writeSyncPlist(path, grapa, operaitionsProto(syncNum))
	return
}

func (p *syncProxy) writeSyncPlist(path string, grapa, plist []byte) (err error) {
	err = p.fs.WriteFileAll(path, plist)
	if err != nil {
		return err
	}
	if len(grapa) == 0 {
		return
	}
	cig := athservice.CigCalc(grapa, plist)
	if len(cig) == 0 {
		return fmt.Errorf("Calc cig fail")
	}
	err = p.fs.WriteFileAll(path+".cig", cig)
	return
}
