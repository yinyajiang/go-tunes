package main

import (
	"fmt"

	itunes "github.com/yinyajiang/go-tunes/iTunesEnv"
)

func main() {
	itu := itunes.New()
	itu.Install(func(Phase string, Prog float64) {
		fmt.Println(Phase, ":", Prog)
	})
}
