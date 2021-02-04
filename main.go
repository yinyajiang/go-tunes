package main

import (
	"fmt"

	iapi "github.com/yinyajiang/go-tunes/iTunesApi"
	mtunes "github.com/yinyajiang/go-tunes/mTunes"
)

func main() {
	mb := iapi.MakeCFBool(true)
	data := iapi.CFToPlist(mb)

	var b bool
	format, err := mtunes.UnmashalPlist(data, &b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format, " ", b)

	pl, err := mtunes.MashalPlist(false)
	if err != nil {
		fmt.Println(err)
	}
	cfb := iapi.PlistToCF(pl)
	fmt.Println(iapi.CFBoolToBool(cfb))
}
