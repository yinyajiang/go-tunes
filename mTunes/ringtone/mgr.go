package ringtone

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/athservice"
	"github.com/yinyajiang/go-tunes/mTunes/fileservice"
	tools "github.com/yinyajiang/go-ytools/utils"
)

//ManagerImpl ...
type ManagerImpl struct {
	dev       mtunes.IOSDevice
	lastDBUTC int64
	tracks    map[uint64]*TrackInfo
}

//New ...
func New(dev mtunes.IOSDevice) (mgr Manager, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device not trusted")
		return
	}
	mgr = &ManagerImpl{
		dev:    dev,
		tracks: make(map[uint64]*TrackInfo, 1),
	}
	return
}

//LoadTrack ...
func (m *ManagerImpl) LoadTrack() (ret map[uint64]*TrackInfo, err error) {
	fs, err := fileservice.New(m.dev)
	if err != nil {
		return
	}
	finfo := fs.GetFileInfo("/iTunes_Control/iTunes/Ringtones.plist")
	if finfo.Modify == m.lastDBUTC {
		ret = m.tracks
		return
	}
	pl, err := fs.ReadFileAll("/iTunes_Control/iTunes/Ringtones.plist")
	if err != nil {
		return
	}
	m.tracks = Parse(pl)
	for _, t := range m.tracks {
		finfo := fs.GetFileInfo(t.Path)
		t.Size = finfo.Size
	}
	m.lastDBUTC = finfo.Modify
	ret = m.tracks
	return
}

//ImportTrack ...
func (m *ManagerImpl) ImportTrack(base ImportTrackInfo) {

	genFileName := func() string {
		name := ""
		for len(name) == 0 {
			name = mtunes.RadomName(4, true)
			for _, t := range m.tracks {
				if t.Name == name && !t.isDeleted {
					name = ""
					break
				}
			}
		}
		return name
	}

	genName := func(bn string) string {
		name := bn
		for _, t := range m.tracks {
			if t.Name == bn && !t.isDeleted {
				name = ""
				break
			}
		}

		i := 1
		for len(name) == 0 {
			name = fmt.Sprintf(bn+"(%d)", i)
			for _, t := range m.tracks {
				if t.Name == name && !t.isDeleted {
					name = ""
					break
				}
			}
			i++
		}
		return name
	}

	var info TrackInfo
	info.FileName = genFileName() + path.Ext(base.SrcPath)
	info.Name = genName(base.Name)
	info.GUID = gen16bitGUID()
	info.TotalTime = base.TotalTime
	info.PID = uint64(tools.RandNum())
	info.Path = "/iTunes_Control/Ringtones/" + info.FileName
	info.Size = base.Size

	info.dev = base.Dev
	info.isFake = true
	info.fakeSrcPath = base.SrcPath

	m.tracks[info.PID] = &info
	return
}

//DeleteTrack ...
func (m *ManagerImpl) DeleteTrack(pid uint64) {
	track, ok := m.tracks[pid]
	if ok {
		track.isDeleted = true
	}
	return
}

//Commit ...
func (m *ManagerImpl) Commit(ctx context.Context) (err error) {
	for pid, track := range m.tracks {
		if track.isFake && track.isDeleted {
			delete(m.tracks, pid)
		}
	}
	inserts := make(map[uint64]*TrackInfo, 0)
	dels := make(map[uint64]*TrackInfo, 0)
	for pid, track := range m.tracks {
		if track.isFake {
			inserts[pid] = track
		} else if track.isDeleted {
			dels[pid] = track
		}
	}
	if len(inserts) == 0 && len(dels) == 0 {
		return
	}
	athProxy := &RingtoneAthProxy{
		inserts: inserts,
		dels:    dels,
		dev:     m.dev,
	}
	ath, err := athservice.New(m.dev, athProxy)
	if err != nil {
		return
	}
	err = ath.Dial()
	if err != nil {
		return
	}
	err = ath.Exec(ctx)
	return
}

func gen16bitGUID() string {
	var b [8]byte
	// Timestamp, 4 bytes
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Pid 2 bytes
	pid := os.Getpid()
	b[4] = byte(pid >> 8)
	b[5] = byte(pid)
	//randnum, 2 bytes, big endian
	num := tools.RandNumN(255)
	b[6] = byte(num >> 8)
	b[7] = byte(num)
	return strings.ToUpper(hex.EncodeToString(b[:]))
}
