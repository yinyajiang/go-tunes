package main

import (
	"context"
	"fmt"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
	"github.com/yinyajiang/go-tunes/mTunes/ringtone"
)

func main() {
	dev := mtunes.OnceForWaitDevice(context.Background(), "2297112500076590" /*"2873943352410158"*/)
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
	fmt.Println(len(tracks))

	for _, t := range tracks {
		mgr.DeleteTrack(t.PID)
	}

	mgr.Commit(context.Background())

	tracks, err = mgr.LoadTrack()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(tracks))
	fmt.Println("Finish")

}
