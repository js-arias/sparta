// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build windows

package win32

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/AllenDang/w32"
	"github.com/js-arias/sparta"
)

// WidgetTable holds a list of widgets.
var widgetTable = make(map[w32.HWND]sparta.Widget)

var (
	// Window extra borders
	extraX = (2 * w32.GetSystemMetrics(w32.SM_CXFRAME)) + 4
	extraY = (2 * w32.GetSystemMetrics(w32.SM_CYFRAME)) + w32.GetSystemMetrics(w32.SM_CYCAPTION) + 4
)

// Window holds the window information.
type window struct {
	id  w32.HWND // window id
	w   sparta.Widget
	pos image.Point

	// graphic part
	dc         w32.HDC // device context
	isPaint    bool    // true if the window is processing a PAINT event.
	back, fore *brush
	curr       *brush
}

func init() {
	sparta.NewWindow = newWindow
}

// NewWindow creates a new window and assigns it to a widget.
func newWindow(w sparta.Widget) {
	var win *window
	rect := w.Property(sparta.Geometry).(image.Rectangle)
	if p := w.Property(sparta.Parent); p != nil {
		pW := p.(sparta.Widget)
		pWin := pW.Window().(*window)
		win = &window{
			w:    w,
			back: pWin.back,
			fore: pWin.fore,
			pos:  rect.Min,
		}
		count := len(pW.Property(sparta.Childs).([]sparta.Widget))
		pW.SetProperty(sparta.Childs, w)
		win.id = w32.CreateWindowEx(0, stringToUTF16(childClass), nil,
			uint(w32.WS_CHILDWINDOW|w32.WS_VISIBLE),
			rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			pWin.id, w32.HMENU(count),
			w32.HINSTANCE(w32.GetWindowLong(pWin.id, w32.GWL_HINSTANCE)), nil)
		if win.id == 0 {
			log.Printf("w32: error: %v\n", getLastError())
			os.Exit(1)
		}
	} else {
		win = &window{
			w:    w,
			back: bkGround,
			fore: frGround,
		}
		win.id = w32.CreateWindowEx(uint(w32.WS_EX_CLIENTEDGE),
			stringToUTF16(baseClass), stringToUTF16(""),
			uint(w32.WS_OVERLAPPEDWINDOW),
			150, 150, rect.Dx()+extraX, rect.Dy()+extraY,
			0, 0, instance, nil)
		if win.id == 0 {
			log.Printf("w32: error: %v\n", getLastError())
			os.Exit(1)
		}
	}
	widgetTable[win.id] = w
	w.SetWindow(win)

	w32.ShowWindow(win.id, w32.SW_SHOWDEFAULT)
	if !w32.UpdateWindow(win.id) {
		log.Printf("w32: error: %v\n", getLastError())
		os.Exit(1)
	}
}

// Close closes the window.
func (win *window) Close() {
	// close the childs
	vc := win.w.Property(sparta.Childs)
	if vc != nil {
		for _, c := range vc.([]sparta.Widget) {
			c.Window().Close()
		}
	}
	delete(widgetTable, win.id)

	if win.w.Property(sparta.Parent) != nil {
		win.w.SetProperty(sparta.Parent, nil)
	}
	win.w.SetProperty(sparta.Childs, nil)
	win.w.RemoveWindow()
	win.w = nil

	w32.DestroyWindow(win.id)

	// if there are no more windows, close the app
	if len(widgetTable) == 0 {
		closeApp()
	}
}

// SetProperty sets a window property.
func (win *window) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Caption:
		if win.w.Property(sparta.Parent) != nil {
			break
		}
		w32.SetWindowText(win.id, v.(string))
	case sparta.Geometry:
		val := v.(image.Rectangle)
		w32.MoveWindow(win.id, val.Min.X, val.Min.Y, val.Dx(), val.Dy(), true)
	case sparta.Foreground:
		val := v.(color.RGBA)
		win.fore = getBrush(val)
	case sparta.Background:
		val := v.(color.RGBA)
		win.back = getBrush(val)
	}
}

// Update updates the window content.
func (win *window) Update() {
	w32.InvalidateRect(win.id, nil, true)
}

// Focus set the focus on the window.
func (win *window) Focus() {
	w32.SetFocus(win.id)
}
