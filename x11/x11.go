// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build darwin freebsd linux netbsd openbsd

// Package x11 defines the x-window backend for sparta.
package x11

import (
	"image/color"
	"log"
	"os"
	"runtime"

	"github.com/js-arias/sparta"
	"github.com/js-arias/xgb"
)

func init() {
	runtime.LockOSThread()
}

var xwin *xgb.Conn

var (
	// WmDeleteWindow event data.
	wmProtocols xgb.Id
	atomType    xgb.Id
	atomDel     xgb.Id
	atomMsg     xgb.Id
	wmDelete    []byte

	// list of allocated pixels
	pixelMap = make(map[uint32]uint32)

	keysyms [256][]int
)

// x11 fixed font
const fixed = "-misc-fixed-medium-r-semicondensed--13-120-75-75-c-60-iso8859-1"

func init() {
	var err error
	xwin, err = xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		log.Printf("x11: error: cannot connect: %v\n", err)
		os.Exit(1)
	}

	// Prepare the WmDeleteWindow event
	protName := "WM_PROTOCOLS"
	wmProt, _ := xwin.InternAtom(false, protName)
	wmProtocols = wmProt.Atom
	atomName := "ATOM"
	atomTp, _ := xwin.InternAtom(false, atomName)
	atomType = atomTp.Atom
	wmDel := "WM_DELETE_WINDOW"
	atmDel, _ := xwin.InternAtom(false, wmDel)
	atomDel = atmDel.Atom
	wmDelete = make([]byte, 4)
	wmDelete[0] = byte(atomDel)
	wmDelete[1] = byte(atomDel >> 8)
	wmDelete[2] = byte(atomDel >> 16)
	wmDelete[3] = byte(atomDel >> 32)

	// Set intern messages
	atmMsg, _ := xwin.InternAtom(false, "SPARTAMSG")
	atomMsg = atmMsg.Atom

	// Prepare pixel maps
	s := xwin.DefaultScreen()
	pixelMap[getColorCode(color.RGBA{R: 255, G: 255, B: 255})] = s.WhitePixel
	pixelMap[getColorCode(color.RGBA{})] = s.BlackPixel

	setKeyboard()

}

func init() {
	sparta.WidthUnit = 6
	sparta.HeightUnit = 13 + 2
}

func getColorCode(c color.RGBA) uint32 {
	code := uint32(c.R) | (uint32(c.G) << 8) | (uint32(c.B) << 16)
	return code
}

func setKeyboard() {
	kmap, _ := xwin.GetKeyboardMapping(xwin.Setup.MinKeycode, xwin.Setup.MaxKeycode-xwin.Setup.MinKeycode+1)
	b := make([]int, 256*int(kmap.KeysymsPerKeycode))
	for i, _ := range keysyms {
		keysyms[i] = b[i*int(kmap.KeysymsPerKeycode) : (i+1)*int(kmap.KeysymsPerKeycode)]
	}
	for i := int(xwin.Setup.MinKeycode); i <= int(xwin.Setup.MaxKeycode); i++ {
		j := i - int(xwin.Setup.MinKeycode)
		for x := 0; x < int(kmap.KeysymsPerKeycode); x++ {
			keysyms[i][x] = int(kmap.Keysyms[(j*int(kmap.KeysymsPerKeycode))+x])
		}
	}
}
