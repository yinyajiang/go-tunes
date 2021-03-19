package fileservice

import (
	"fmt"
	"io"

	mtunes "github.com/yinyajiang/go-tunes"
	iapi "github.com/yinyajiang/go-tunes/itunesapi"
)

type fileImpl struct {
	hand       uint64
	afcConnect uintptr
	info       *FileInfo
	dev        mtunes.Device
}

//Read ...
func (f *fileImpl) Read(p []byte) (readed int, err error) {
	if f.hand == 0 {
		err = fmt.Errorf("File hand is 0")
		return
	}
	for {
		select {
		case <-f.dev.ExtrackContext().Done():
			err = fmt.Errorf("Device is extrack")
			return
		default:
		}

		read := len(p) - readed
		if read <= 0 {
			break
		}

		if read > 1024*1024*4 {
			read = 1024 * 1024 * 4
		}

		oneReaded, res := iapi.AFCFileRefRead(f.afcConnect, f.hand, p[readed:readed+read])
		if res != 0 {
			err = fmt.Errorf("Afc Read fail")
			break
		}
		if oneReaded > read {
			err = fmt.Errorf("AFCFileRefRead readed > read")
			break
		}
		if oneReaded == 0 {
			err = io.EOF
			break
		}
		readed += oneReaded
	}
	return
}

//Close ...
func (f *fileImpl) Close() (err error) {
	if f.hand != 0 {
		iapi.AFCFileRefClose(f.afcConnect, f.hand)
		f.hand = 0
	}
	return nil
}

//Write ...
func (f *fileImpl) Write(p []byte) (writen int, err error) {
	if f.hand == 0 {
		err = fmt.Errorf("File hand is 0")
		return
	}
	for {
		select {
		case <-f.dev.ExtrackContext().Done():
			err = fmt.Errorf("Device is extrack")
			return
		default:
		}

		write := len(p) - writen
		if write <= 0 {
			break
		}
		if write > 1024*1024*4 {
			write = 1024 * 1024 * 4
		}
		res := iapi.AFCFileRefWrite(f.afcConnect, f.hand, p[writen:writen+write])
		if res != 0 {
			err = fmt.Errorf("Afc Write fail")
			break
		}
		writen += write
	}
	return
}

//Seek ...
func (f *fileImpl) Seek(offset int64, whence int) (curSeek int64, err error) {
	if f.hand == 0 {
		err = fmt.Errorf("File hand is 0")
		return
	}
	if 0 != iapi.AFCFileRefSeek(f.afcConnect, f.hand, uint64(offset), whence) {
		err = fmt.Errorf("AFC Seek fail")
		return
	}

	if 0 == whence {
		curSeek = offset
	} else {
		t, res := iapi.AFCFileRefTell(f.afcConnect, f.hand)
		if res != 0 {
			err = fmt.Errorf("AFC Seek:Tell fail")
			return
		}
		curSeek = int64(t)
	}
	return
}

//ReadAt ...
func (f *fileImpl) ReadAt(p []byte, off int64) (n int, err error) {
	if f.hand == 0 {
		err = fmt.Errorf("File hand is 0")
		return
	}
	_, err = f.Seek(off, 0)
	if err != nil {
		return
	}
	return f.Read(p)
}

//WriteAt ...
func (f *fileImpl) WriteAt(p []byte, off int64) (n int, err error) {
	if f.hand == 0 {
		err = fmt.Errorf("File hand is 0")
		return
	}
	_, err = f.Seek(off, 0)
	if err != nil {
		return
	}
	return f.Write(p)
}
