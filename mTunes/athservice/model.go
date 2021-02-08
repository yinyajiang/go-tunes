package athservice

import "io"

//AirProxy ...
type AirProxy interface {
	SubmitReadyPlist(syncNum string, grapa []byte) error
	OpenEntityReader(entityID string) (r io.ReadCloser, size int64, err error)
	OpenEntityWriter(entityID string) (w io.WriteCloser, target string, err error)
	TransferFinish(entityID string)
	GetKeybag() (keybag string, anchors []string)
}
