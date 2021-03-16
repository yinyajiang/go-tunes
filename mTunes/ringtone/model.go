package ringtone

import (
	"context"

	mtunes "github.com/yinyajiang/go-tunes/mtunes"
)

//TrackInfo ...
type TrackInfo struct {
	//base
	FileName  string `json:"fileName"`
	Name      string `json:"name"`
	GUID      string `json:"guid"`
	TotalTime uint64 `json:"totalTime"`
	PID       uint64 `json:"pid"`
	Path      string `json:"path"`
	Protected bool   `json:"protected"`
	Size      int64  `json:"size"`
	Purchased bool   `json:"purchased"`

	isDeleted bool
	isInvalid bool

	//fake
	isFake      bool
	fakeSrcPath string
	dev         mtunes.Device
}

//ImportTrackInfo ...
type ImportTrackInfo struct {
	Name      string
	TotalTime uint64
	Size      int64
	Dev       mtunes.Device
	SrcPath   string
}

//Manager ...
type Manager interface {
	LoadTrack() (ret []TrackInfo, err error)
	ImportTrack(base ImportTrackInfo)
	DeleteTrack(pid uint64)
	Commit(ctx context.Context) error
}
