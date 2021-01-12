package main

/*
#cgo CXXFLAGS: -std=c++11
#include "tunesApi.h"
*/
import "C"

func main() {
	C.AFCConnectionClose(nil)
}
