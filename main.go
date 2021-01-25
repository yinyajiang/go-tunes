package main

import (
	"context"
	"fmt"
	"time"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
)

func main() {
	mtunes.SubscriptionForWaitOperatior()
	go func() {
		con, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
		dev := mtunes.WaitForDevice(con, "2873943352410158")
		if dev == nil {
			fmt.Println("time out")
			mtunes.StopEventLoop()
			return
		}
		con2, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
		err := dev.WaitTrust(con2)
		if nil == err {
			fmt.Println("trust success")
		} else {
			fmt.Println(err)
		}
		fmt.Println(dev.DeviceInfo())
	}()
	mtunes.RunEventLoop()
	fmt.Println("finish")

}
