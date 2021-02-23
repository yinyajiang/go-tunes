package main

import (
	"context"
	"fmt"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/ringtone"
)

func main() {
	dev := mtunes.OnceForWaitDevice(context.Background(), "2873943352410158")
	if dev == nil {
		fmt.Println("Not found dev")
		return
	}
	mgr, err := ringtone.New(dev)
	if err != nil {
		fmt.Println(err)
		return
	}

	tracks, err := mgr.LoadTrack()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range tracks {
		fmt.Println(t)
	}

	base := ringtone.ImportTrackInfo{
		Name:      "没诶13",
		TotalTime: 242284,
		Size:      10780022,
		SrcPath:   "/Users/new/Documents/Temp/贰佰 - 玫瑰 [mqms2] [高质量].m4r",
	}
	mgr.ImportTrack(base)
	mgr.Commit(context.Background())

	tracks, err = mgr.LoadTrack()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range tracks {
		fmt.Println(t)
	}
	fmt.Println("Finish")

}
