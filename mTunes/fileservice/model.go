package fileservice

import "io"

const (
	IFT_Dir = iota
	IFT_Link
	IFT_File
	IFT_Blk
	IFT_Chr
	IFT_Fifo
	IFT_Mt
	IFT_Sock
	IFT_Unknown
)

const (
	AFC_FOPEN_RDONLY   = 0x00000001
	AFC_FOPEN_RW       = 0x00000002
	AFC_FOPEN_WRONLY   = 0x00000003
	AFC_FOPEN_WRT      = 0x00000004
	AFC_FOPEN_APPEND   = 0x00000005
	AFC_FOPEN_RDAPPEND = 0x00000006
)

const (
	AFC_SEEK_SET = 0
	AFC_SEEK_CUR = 1
	AFC_SEEK_END = 2
)

//FileInfo ...
type FileInfo struct {
	Name   string
	Size   int64
	Create int64
	Modify int64
	Type   int
}

//File ...
type File interface {
	io.ReadWriteSeeker
	io.Closer
	io.ReaderAt
	io.WriterAt
}

//FileService ...
type FileService interface {
	GetFileInfo(path string) *FileInfo
	IsFileExist(path string) bool
	PathWalk(dir string, dirFun func(path string, info *FileInfo, postName string)) bool
	CreateDirectorys(path string)
	RemovePath(path string)
	OpenFile(path string, mode int64) File
}
