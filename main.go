package main

import (
	"context"
	"fmt"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/ringtone"
)

func main() {
	dev := mtunes.OnceEventLoopForWaitDevice(context.Background(), "2873943352410158")
	if dev == nil {
		fmt.Println("Not found dev")
		return
	}
	mgr, err := ringtone.New(dev)
	if err != nil {
		fmt.Println(err)
		return
	}

	base := ringtone.ImportTrackInfo{
		Name:      "case0_山丹丹开花红艳艳",
		TotalTime: 242284,
		Size:      10780022,
		SrcPath:   "/Volumes/192.168.0.85/test file/case0_山丹丹开花红艳艳.mp3",
	}
	mgr.ImportTrack(base)
	mgr.Commit(context.Background())

	tracks, err := mgr.LoadTrack()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range tracks {
		fmt.Println(*t)
	}

}
