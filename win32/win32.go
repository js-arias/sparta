// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build windows

// Package win32 defines the windows backend for sparta.
package win32

import (
	"image/color"
	"log"
	"os"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/AllenDang/w32"
	"github.com/js-arias/sparta"
)

func init() {
	runtime.LockOSThread()
}

const (
	// Window classes
	baseClass  = "BaseSpartaClass"
	childClass = "ChildSpartaClass"
)

// global handle
var instance = w32.GetModuleHandle("")

func init() {
	// base class
	wc := &w32.WNDCLASSEX{
		Style:      w32.CS_HREDRAW | w32.CS_VREDRAW,
		WndProc:    syscall.NewCallback(winEvent),
		Instance:   instance,
		Icon:       w32.LoadIcon(0, w32.MakeIntResource(w32.IDI_APPLICATION)),
		Cursor:     w32.LoadCursor(0, w32.MakeIntResource(w32.IDC_ARROW)),
		Background: createSolidBrush(rgb(color.RGBA{R: 255, G: 255, B: 255})),
		ClassName:  stringToUTF16(baseClass),
	}
	wc.IconSm = wc.Icon
	wc.Size = uint32(unsafe.Sizeof(*wc))
	if w32.RegisterClassEx(wc) == 0 {
		log.Printf("w32: error: %v\n", getLastError())
		os.Exit(1)
	}

	// child class
	wc.WndExtra = 32
	wc.Icon, wc.IconSm = 0, 0
	wc.ClassName = stringToUTF16(childClass)
	w32.RegisterClassEx(wc)
}

var (
	// general objetcs
	winFont  w32.HFONT
	bkGround *brush
	frGround *brush
	mapBrush = make(map[w32.COLORREF]*brush)
)

func init() {
	font := &w32.LOGFONT{}
	font.Height = 12
	fontName := utf16.Encode([]rune("Lucida Console"))
	for i, r := range fontName {
		font.FaceName[i] = r
	}
	winFont = w32.CreateFontIndirect(font)
	frGround = getBrush(color.RGBA{})
	bkGround = getBrush(color.RGBA{R: 255, G: 255, B: 255})
}

func init() {
	sparta.WidthUnit = 7
	sparta.HeightUnit = 12 + 2
}
