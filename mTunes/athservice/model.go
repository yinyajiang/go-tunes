package athservice

import (
	"context"
	"io"
)

//AthProxy ...
type AthProxy interface {
	SubmitReadyPlist(syncNum string, grapa []byte) error
	IsExistAsset(entityID string) bool
	OpenAssetReader(entityID string) (r io.ReadCloser, size int64, err error)
	OpenAssetWriter(entityID string) (w io.WriteCloser, target string, err error)
	AssetFinish(entityID string, successed bool)
	GetKeybag() (keybag string, anchors []string)
}

//Service ...
type Service interface {
	Dial() (err error)
	Serve(ctx context.Context) (err error)
}
