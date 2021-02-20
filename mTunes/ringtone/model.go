package ringtone

import (
	"context"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
)

//TrackInfo ...
type TrackInfo struct {
	//base
	FileName  string
	Name      string
	GUID      string
	TotalTime uint64
	PID       uint64
	Path      string
	Protected bool
	Size      int64

	isDeleted bool

	//fake
	isFake      bool
	fakeSrcPath string
	dev         mtunes.IOSDevice
}

//ImportTrackInfo ...
type ImportTrackInfo struct {
	Name      string
	TotalTime uint64
	Size      int64
	Dev       mtunes.IOSDevice
	SrcPath   string
}

//Manager ...
type Manager interface {
	LoadTrack() (ret map[uint64]*TrackInfo, err error)
	ImportTrack(base ImportTrackInfo)
	DeleteTrack(pid uint64)
	Commit(ctx context.Context) error
}
