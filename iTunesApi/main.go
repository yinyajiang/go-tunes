package main

import "C"
import (
	"fmt"
	"time"
)

func main() {
	AMRestorableDeviceRegisterForNotificationsForDevices(func(even, mode string, modeDevice, restorableDevice uintptr) {
		fmt.Println(even, " ", mode, " ", modeDevice)
	})
	go func() {
		select {
		case <-time.After(time.Second * 10):
			RunLoopStop()
		}
	}()
	RunLoopRun()
	fmt.Println("finish")
}
