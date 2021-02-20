package ringtone

import (
	"fmt"
	"io"
	"strconv"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/athservice"
	"github.com/yinyajiang/go-tunes/mTunes/fileservice"
	tools "github.com/yinyajiang/go-ytools/utils"
)

//RingtoneAthProxy ...
type RingtoneAthProxy struct {
	inserts map[uint64]*TrackInfo
	dels    map[uint64]*TrackInfo
	dev     mtunes.IOSDevice
}

func (p *RingtoneAthProxy) TransferFinish(entityID string) {

}

//GetKeybag ...
func (p *RingtoneAthProxy) GetKeybag() (keybag string, anchors []string) {
	return "Ringtone", []string{"0"}
}

//OpenEntityWriter ...
func (p *RingtoneAthProxy) OpenEntityWriter(entityID string) (w io.WriteCloser, target string, err error) {
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
	fs, err := fileservice.New(track.dev)
	w, err = fs.OpenFile(track.Path, fileservice.AFC_FOPEN_WRONLY)
	target = track.Path
	return
}

//OpenEntityReader ...
func (p *RingtoneAthProxy) OpenEntityReader(entityID string) (r io.ReadCloser, size int64, err error) {
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
		var fs fileservice.FileService
		fs, err = fileservice.New(track.dev)
		r, err = fs.OpenFile(track.fakeSrcPath, fileservice.AFC_FOPEN_RDONLY)
	}

	return
}

//SubmitReadyPlist ...
func (p *RingtoneAthProxy) SubmitReadyPlist(syncNum string, grapa []byte) (err error) {
	insertArr := make([]*TrackInfo, 0, len(p.inserts))
	delArr := make([]*TrackInfo, 0, len(p.dels))
	for _, t := range p.inserts {
		insertArr = append(insertArr, t)
	}
	for _, t := range p.dels {
		delArr = append(delArr, t)
	}

	n, _ := strconv.Atoi(syncNum)
	pl, err := mashalUpdateProto(n, insertArr, delArr)
	if err != nil {
		return
	}
	path := fmt.Sprintf("/iTunes_Control/Ringtones/Sync/Sync_%08d.plist", uint(n))
	if err != nil {
		return
	}
	err = p.writeSyncPlist(path, nil, pl)
	if err != nil {
		return
	}
	path = fmt.Sprintf("/iTunes_Control/Sync/Media/Sync_%08d.plist", uint(n))
	if err != nil {
		return
	}
	err = p.writeSyncPlist(path, grapa, operaitionsProto())
	return
}

func (p *RingtoneAthProxy) writeSyncPlist(path string, grapa, plist []byte) (err error) {
	fs, err := fileservice.New(p.dev)
	if err != nil {
		return
	}
	err = fs.WriteFileAll(path, plist)
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
	err = fs.WriteFileAll(path+".cig", cig)
	return
}
