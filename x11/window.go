// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build darwin freebsd linux netbsd openbsd

package x11

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
	"github.com/js-arias/xgb"
)

// widgetTable holds a list of widgets.
var widgetTable = make(map[xgb.Id]sparta.Widget)

// Window holds the window information.
type window struct {
	id xgb.Id
	w  sparta.Widget // associated widget

	// graphic part
	gc         xgb.Id // graphic context
	isExpose   bool   //true if the window is processing an expose event.
	back, fore uint32
}

func init() {
	sparta.NewWindow = newWindow
}

// NewWindow creates a new window and assigns it to a widget.
func newWindow(w sparta.Widget) {
	s := xwin.DefaultScreen()
	pId := s.Root
	win := &window{
		id:   xwin.NewId(),
		w:    w,
		gc:   xwin.NewId(),
		back: s.WhitePixel,
		fore: s.BlackPixel,
	}
	if p := w.Property(sparta.Parent); p != nil {
		pw := p.(sparta.Widget)
		pWin := pw.Window().(*window)
		win.back = pWin.back
		win.fore = pWin.fore
		pId = pWin.id
		pw.SetProperty(sparta.Childs, w)
	}
	widgetTable[win.id] = w
	w.SetWindow(win)
	r := w.Property(sparta.Geometry).(image.Rectangle)
	xwin.CreateWindow(0, win.id, pId,
		int16(r.Min.X), int16(r.Min.Y), uint16(r.Dx()), uint16(r.Dy()), 0,
		xgb.WindowClassInputOutput, s.RootVisual, 0, nil)
	xwin.ChangeWindowAttributes(win.id, xgb.CWBackPixel|xgb.CWEventMask,
		[]uint32{
			win.back,
			allEventMask,
		})
	font := xwin.NewId()
	xwin.OpenFont(font, fixed)
	xwin.CreateGC(win.gc, win.id, xgb.GCBackground|xgb.GCForeground|xgb.GCFont,
		[]uint32{
			win.fore,
			win.back,
			uint32(font),
		})
	xwin.CloseFont(font)
	xwin.MapWindow(win.id)
	xwin.ChangeProperty(xgb.PropModeReplace, win.id, wmProtocols, atomType, 32, wmDelete)
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

	xwin.DestroyWindow(win.id)

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
		xwin.ChangeProperty(xgb.PropModeReplace, win.id, xgb.AtomWmName,
			xgb.AtomString, 8, []byte(v.(string)))
	case sparta.Geometry:
		val := v.(image.Rectangle)
		xwin.ConfigureWindow(win.id, xgb.ConfigWindowX|xgb.ConfigWindowY|xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
			[]uint32{
				uint32(val.Min.X),
				uint32(val.Min.Y),
				uint32(val.Dx()),
				uint32(val.Dy()),
			})
	case sparta.Foreground:
		val := v.(color.RGBA)
		s := xwin.DefaultScreen()
		code := getColorCode(val)
		px, ok := pixelMap[code]
		if !ok {
			r, g, b, _ := val.RGBA()
			cl, _ := xwin.AllocColor(s.DefaultColormap, uint16(r), uint16(g), uint16(b))
			px = cl.Pixel
			pixelMap[code] = px
		}
		win.fore = px
	case sparta.Background:
		val := v.(color.RGBA)
		s := xwin.DefaultScreen()
		code := getColorCode(val)
		px, ok := pixelMap[code]
		if !ok {
			r, g, b, _ := val.RGBA()
			cl, _ := xwin.AllocColor(s.DefaultColormap, uint16(r), uint16(g), uint16(b))
			px = cl.Pixel
			pixelMap[code] = px
		}
		win.back = px
	}
}

// Update updates thw window.
func (win *window) Update() {
	rect := win.w.Property(sparta.Geometry).(image.Rectangle)
	xwin.ClearArea(true, win.id, 0, 0, uint16(rect.Dx()), uint16(rect.Dy()))
}

// Focus set the focus on the window.
func (win *window) Focus() {
	xwin.SetInputFocus(xgb.InputFocusNone, win.id, xgb.TimeCurrentTime)
}
