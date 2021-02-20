package athservice

import (
	"context"
	"io"
)

//AthProxy ...
type AthProxy interface {
	SubmitReadyPlist(syncNum string, grapa []byte) error
	OpenEntityReader(entityID string) (r io.ReadCloser, size int64, err error)
	OpenEntityWriter(entityID string) (w io.WriteCloser, target string, err error)
	TransferFinish(entityID string)
	GetKeybag() (keybag string, anchors []string)
}

//AthService ...
type AthService interface {
	Dial() (err error)
	Exec(ctx context.Context) (err error)
}
