package tunes

import "context"

//Device ...
type Device interface {
	Mode() string
	ID() string
	ECID() string
	Model() string
	Version() string
	ModelName() string
	DeviceInfo() map[string]interface{}
	ModeDevice() uintptr
	OriginalDevice() uintptr
	Trust() error
	WaitTrust(ctx context.Context) error
	AbordTrust()
	IsTrusted() bool
	IsExtract() bool

	ExtrackContext() context.Context

	SaveUserData(key string, val interface{})
	GetUserData(key string) interface{}
	DeleteUserData(key string)

	GetStartService(name string) (conn uintptr, err error)
	IsServiceRuning(name string) bool
	StopService(name string)

	WorkDir(join string) string

	Release()
}
