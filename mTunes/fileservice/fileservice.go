package fileservice

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	tools "github.com/yinyajiang/go-ytools/utils"
)

//Service ...
type Service struct {
	dev        mtunes.IOSDevice
	afcConnect uintptr
}

//New ...
func New(dev mtunes.IOSDevice) (svc FileService, err error) {
	if dev == nil {
		err = fmt.Errorf("Dev is nil")
		return
	}
	afcConnect, err := getAfcConnect(dev)
	if err != nil {
		return
	}
	psvc := &Service{
		dev:        dev,
		afcConnect: afcConnect,
	}
	svc = psvc
	return
}

//GetFileInfo ...
func (svc *Service) GetFileInfo(path string) *FileInfo {
	hand := iapi.AFCFileInfoOpen(svc.afcConnect, path)
	if 0 == hand {
		return nil
	}
	defer iapi.AFCKeyValueClose(hand)
	var info FileInfo
	info.Name = tools.PathName(path)
	for {
		key, val := iapi.AFCKeyValueRead(hand)
		if len(key) == 0 {
			break
		}

		if strings.EqualFold(key, "st_size") {
			info.Size, _ = strconv.ParseInt(val, 0, 64)
		} else if strings.EqualFold(key, "st_birthtime") {
			info.Create, _ = strconv.ParseInt(val, 0, 64)
		} else if strings.EqualFold(key, "st_mtime") {
			info.Modify, _ = strconv.ParseInt(val, 0, 64)
		} else if strings.EqualFold(key, "st_ifmt") {
			if strings.EqualFold(key, "S_IFDIR") {
				info.Type = IFT_Dir
			} else if strings.EqualFold(key, "S_IFLNK") {
				info.Type = IFT_Link
			} else if strings.EqualFold(key, "S_IFREG") {
				info.Type = IFT_File
			} else if strings.EqualFold(key, "S_IFBLK") {
				info.Type = IFT_Blk
			} else if strings.EqualFold(key, "S_IFCHR") {
				info.Type = IFT_Chr
			} else if strings.EqualFold(key, "S_IFIFO") {
				info.Type = IFT_Fifo
			} else if strings.EqualFold(key, "S_IFMT") {
				info.Type = IFT_Mt
			} else if strings.EqualFold(key, "S_IFSOCK") {
				info.Type = IFT_Sock
			} else {
				info.Type = IFT_Unknown
			}
		}
	}
	return &info
}

//IsFileExist ...
func (svc *Service) IsFileExist(path string) bool {
	hand := iapi.AFCFileInfoOpen(svc.afcConnect, path)
	if 0 == hand {
		return false
	}
	defer iapi.AFCKeyValueClose(hand)
	return true
}

//PathWalk ...
func (svc *Service) PathWalk(dir string, dirFun func(string, *FileInfo, string) bool) {
	if dirFun == nil {
		return
	}
	hand := iapi.AFCDirectoryOpen(svc.afcConnect, dir)
	if 0 == hand {
		fmt.Println("Open Dir fail")
		return
	}
	defer iapi.AFCDirectoryClose(svc.afcConnect, hand)

	for {
		select {
		case <-svc.dev.ExtrackContext().Done():
			return
		default:
		}

		name := iapi.AFCDirectoryRead(svc.afcConnect, hand)
		if len(name) == 0 {
			break
		}
		if name == "." || name == ".." {
			continue
		}
		path := tools.ThePath(dir, name)
		info := svc.GetFileInfo(path)
		if info == nil {
			continue
		}
		postName := tools.PostPath(path, dir)
		if !dirFun(path, info, postName) {
			break
		}
	}
	return
}

//CreateDirectorys ...
func (svc *Service) CreateDirectorys(path string) {
	for i, j := 0, 0; i != -1; j++ {
		i = strings.Index(path[j:], "/")
		if i == -1 {
			iapi.AFCDirectoryCreate(svc.afcConnect, path)
			break
		}
		j += i
		if i != 0 {
			iapi.AFCDirectoryCreate(svc.afcConnect, path[0:j])
		}
	}
}

//RemovePath ...
func (svc *Service) RemovePath(path string) {
	iapi.AFCRemovePath(svc.afcConnect, path)
}

//RenameAndMove ...
func (svc *Service) RenameAndMove(src, dst string) {
	iapi.AFCRenamePath(svc.afcConnect, src, dst)
}

//OpenFile ...
func (svc *Service) OpenFile(path string, mode int64) (f File, err error) {
	if mode == AFC_FOPEN_WRONLY {
		svc.RemovePath(path)
	}
	fhand := iapi.AFCFileRefOpen(svc.afcConnect, path, mode)
	if 0 == fhand {
		err = fmt.Errorf("Open %s fail\n", path)
		return
	}
	f = &DeviceFileImpl{
		hand:       fhand,
		afcConnect: svc.afcConnect,
		dev:        svc.dev,
	}
	return
}

//ReadFileAll ...
func (svc *Service) ReadFileAll(path string) (data []byte, err error) {
	f, err := svc.OpenFile(path, AFC_FOPEN_RDONLY)
	if err != nil {
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	return
}

//WriteFileAll ...
func (svc *Service) WriteFileAll(path string, data []byte) (err error) {
	svc.CreateDirectorys(tools.AbsPath(path))
	f, err := svc.OpenFile(path, AFC_FOPEN_WRONLY)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

func getAfcConnect(dev mtunes.IOSDevice) (afcConnect uintptr, err error) {
	//Connect AFC
	afcConnect, ok := dev.GetUserData("afc_connect").(uintptr)
	if ok {
		return
	}

	afc, err := dev.GetStartService("com.apple.afc")
	if err != nil {
		return
	}
	afcConnect = iapi.AFCConnectionOpen(afc)
	if afcConnect == 0 {
		err = fmt.Errorf("Open Afc connect fail")
		return
	}
	dev.SaveUserData("afc_connect", afcConnect)
	return
}
