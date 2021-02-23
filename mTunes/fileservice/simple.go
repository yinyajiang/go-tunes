package fileservice

import mtunes "github.com/yinyajiang/go-tunes/mTunes"

//ReadFileAll ...
func ReadFileAll(dev mtunes.Device, path string) (data []byte, err error) {
	svc, err := New(dev)
	if err != nil {
		return
	}
	defer func() {
		svc.Release()
	}()
	return svc.ReadFileAll(path)
}

//WriteFileAll ...
func WriteFileAll(dev mtunes.Device, path string, data []byte) (err error) {
	svc, err := New(dev)
	if err != nil {
		return
	}
	defer func() {
		svc.Release()
	}()
	return svc.WriteFileAll(path, data)
}

//GetFileInfo ...
func GetFileInfo(dev mtunes.Device, path string) *FileInfo {
	svc, err := New(dev)
	if err != nil {
		return nil
	}
	defer func() {
		svc.Release()
	}()
	return svc.GetFileInfo(path)
}
