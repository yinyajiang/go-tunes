package mtunes

import "context"

//IOSDevice ...
type IOSDevice interface {
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

	SaveUserData(key string, val interface{})
	GetUserData(key string) interface{}
	DeleteUserData(key string)

	GetStartService(name string) (conn uintptr, err error)
	IsServiceRuning(name string) bool
	StopService(name string)
}
