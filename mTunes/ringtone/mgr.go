package ringtone

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"path"
	"strings"
	"time"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/athservice"
	"github.com/yinyajiang/go-tunes/mTunes/fileservice"
	tools "github.com/yinyajiang/go-ytools/utils"
)

type managerImpl struct {
	dev         mtunes.Device
	tracks      map[uint64]*TrackInfo
	nameSet     map[string]struct{}
	fileNameSet map[string]struct{}
	bLoaded     bool
}

//New ...
func New(dev mtunes.Device) (mgr Manager, err error) {
	if !dev.IsTrusted() {
		err = fmt.Errorf("Device not trusted")
		return
	}
	mgr = &managerImpl{
		dev: dev,
	}
	return
}

//LoadTrack ...
func (m *managerImpl) LoadTrack() (ret []TrackInfo, err error) {
	fs, err := fileservice.New(m.dev)
	if err != nil {
		return
	}

	defer func() {
		fs.Release()
		if len(m.tracks) == 0 || err != nil {
			return
		}
		ret = make([]TrackInfo, 0, len(m.tracks))
		for _, t := range m.tracks {
			if !t.isDeleted && !t.isInvalid {
				ret = append(ret, *t)
			}
		}
	}()

	m.bLoaded = true
	pl, err := fs.ReadFileAll("/iTunes_Control/iTunes/Ringtones.plist")
	if err != nil {
		return
	}
	m.tracks, m.nameSet, m.fileNameSet = Parse(pl, fs)
	return
}

//ImportTrack ...
func (m *managerImpl) ImportTrack(base ImportTrackInfo) {
	if !m.bLoaded {
		m.LoadTrack()
	}

	genFileName := func() string {
		name := ""
		for len(name) == 0 {
			name = mtunes.RadomName(4, true)
			_, ok := m.fileNameSet[name]
			if !ok {
				m.fileNameSet[name] = struct{}{}
				break
			}
			name = ""
		}
		return name
	}

	genName := func(bn string) string {
		name := bn
		_, ok := m.nameSet[name]
		if !ok {
			m.nameSet[name] = struct{}{}
			return name
		}
		name = ""
		i := 1
		for len(name) == 0 {
			name = fmt.Sprintf(bn+"(%d)", i)
			_, ok := m.nameSet[name]
			if !ok {
				m.nameSet[name] = struct{}{}
				break
			}
			name = ""
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
func (m *managerImpl) DeleteTrack(pid uint64) {
	if !m.bLoaded {
		m.LoadTrack()
	}

	track, ok := m.tracks[pid]
	if ok {
		track.isDeleted = true
	}
	return
}

//Commit ...
func (m *managerImpl) Commit(ctx context.Context) (err error) {
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

	fs, err := fileservice.New(m.dev)
	if err != nil {
		return
	}
	defer func() {
		fs.Release()
	}()
	ath, err := athservice.New(m.dev, &syncProxy{
		inserts: inserts,
		dels:    dels,
		fs:      fs,
	})
	if err != nil {
		return
	}
	err = ath.Dial()
	if err != nil {
		return
	}
	err = ath.Serve(ctx)

	//remove name set
	if err == nil {
		for _, track := range m.tracks {
			if track.isDeleted {
				delete(m.fileNameSet, track.FileName)
				delete(m.nameSet, track.Name)
			}
		}
	}

	m.LoadTrack()
	return
}

func gen16bitGUID() string {
	var b [8]byte
	// Timestamp, 4 bytes
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	//randnum, 4 bytes
	binary.BigEndian.PutUint32(b[4:], uint32(tools.RandNum()))
	return strings.ToUpper(hex.EncodeToString(b[:]))
}
