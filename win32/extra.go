// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build windows

package win32

import (
	"image/color"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/AllenDang/w32"
)

// GetLastError retrieves ge calling thread's last-error value
func getLastError() error {
	return syscall.GetLastError()
}

// stringToUTF16 returns a pointer to a UTF-16 string
func stringToUTF16(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

// toUFT16 returns a pointer to a UTF-16 string
func toUTF16(s string) (uintptr, int) {
	tx := utf16.Encode([]rune(s))
	return uintptr(unsafe.Pointer(&tx[0])), len(tx)
}

// Window macros

// Rgb returns a rgb color
func rgb(c color.RGBA) w32.COLORREF {
	color := uint32(c.R) | (uint32(c.G) << 8) | (uint32(c.B) << 16)
	return w32.COLORREF(color)
}

func loWord(dw uint32) uint16 {
	return uint16(dw)
}

func hiWord(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func getXLParam(lp uintptr) int {
	return int(int16(loWord(uint32(lp))))
}

func getYLParam(lp uintptr) int {
	return int(int16(hiWord(uint32(lp))))
}
func getWheelDeltaWParam(wparam uintptr) int16 {
	return int16(hiWord(uint32(wparam)))
}

var modgdi32 = syscall.NewLazyDLL("gdi32.dll")

var (
	// gdi
	procSaveDC           = modgdi32.NewProc("SaveDC")
	procRestoreDC        = modgdi32.NewProc("RestoreDC")
	procCreatePen        = modgdi32.NewProc("CreatePen")
	procCreateSolidBrush = modgdi32.NewProc("CreateSolidBrush")
	procTextOut          = modgdi32.NewProc("TextOutW")
	procPolyLine         = modgdi32.NewProc("Polyline")
	procArc              = modgdi32.NewProc("Arc")
	procPie              = modgdi32.NewProc("Pie")
	procPolygon          = modgdi32.NewProc("Polygon")
	procSetPixel         = modgdi32.NewProc("SetPixel")
)

func saveDC(hdc w32.HDC) int {
	ret, _, _ := procSaveDC.Call(uintptr(hdc))
	return int(ret)
}

func restoreDC(hdc w32.HDC, idSaved int) {
	procRestoreDC.Call(uintptr(hdc), uintptr(idSaved))
}

func createPen(style, width int, color w32.COLORREF) w32.HPEN {
	ret, _, _ := procCreatePen.Call(uintptr(style), uintptr(width), uintptr(color))
	return w32.HPEN(ret)
}

func createSolidBrush(color w32.COLORREF) w32.HBRUSH {
	ret, _, _ := procCreateSolidBrush.Call(uintptr(color))
	return w32.HBRUSH(ret)
}

func textOut(hdc w32.HDC, xStart, yStart int, str string) {
	tx, sz := toUTF16(str + "\n")
	procTextOut.Call(uintptr(hdc), uintptr(xStart), uintptr(yStart), tx, uintptr(sz))
}

func polyLine(hdc w32.HDC, ppt []w32.POINT) bool {
	ret, _, _ := procPolyLine.Call(uintptr(hdc), uintptr(unsafe.Pointer(&ppt[0])),
		uintptr(len(ppt)))
	return ret != 0
}

func arc(hdc w32.HDC, left, top, right, bottom, xStart, yStart, xEnd, yEnd int) bool {
	ret, _, _ := procArc.Call(uintptr(hdc), uintptr(left), uintptr(top),
		uintptr(right), uintptr(bottom), uintptr(xStart), uintptr(yStart),
		uintptr(xEnd), uintptr(yEnd))
	return ret != 0
}

func pie(hdc w32.HDC, left, top, right, bottom, xStart, yStart, xEnd, yEnd int) bool {
	ret, _, _ := procPie.Call(uintptr(hdc), uintptr(left), uintptr(top),
		uintptr(right), uintptr(bottom), uintptr(xStart), uintptr(yStart),
		uintptr(xEnd), uintptr(yEnd))
	return ret != 0
}

func polygon(hdc w32.HDC, ppt []w32.POINT) bool {
	ret, _, _ := procPolygon.Call(uintptr(hdc), uintptr(unsafe.Pointer(&ppt[0])),
		uintptr(len(ppt)))
	return ret != 0
}

func setPixel(hdc w32.HDC, x, y int, color w32.COLORREF) w32.COLORREF {
	ret, _, _ := procSetPixel.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(color))
	return w32.COLORREF(ret)
}

var moduser32 = syscall.NewLazyDLL("user32.dll")

var procGetKeyState = moduser32.NewProc("GetKeyState")

func GetKeyState(nVirtKey int) int16 {
	ret, _, _ := procGetKeyState.Call(uintptr(nVirtKey))
	return int16(ret)
}
