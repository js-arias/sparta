// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build darwin freebsd linux netbsd openbsd

package x11

import (
	"image"
	"image/color"
	"math"
	"unicode/utf16"

	"github.com/js-arias/sparta"
	"github.com/js-arias/xgb"
)

// Draw sets the drawing mode of the window.
func (win *window) Draw(mode bool) {
	if win.isExpose {
		return
	}
	if mode {
		xwin.ChangeGC(win.gc, xgb.GCForeground, []uint32{win.fore})
	}
}

// Text draws text in the window.
func (win *window) Text(pt image.Point, text string) {
	if len(text) == 0 {
		return
	}
	if len(text) > 128 {
		text = text[:128]
	}
	tx := utf16.Encode([]rune(text))
	tx16 := make([]xgb.Char2b, len(tx))
	for i, v := range tx {
		tx16[i].Byte1 = byte(v >> 8)
		tx16[i].Byte2 = byte(v)
	}
	x16, y16 := int16(pt.X), int16(pt.Y+9)
	xwin.ImageText16(win.id, win.gc, x16, y16, tx16)
}

// Rectangle draws a rectangle in the window.
func (win *window) Rectangle(rect image.Rectangle, fill bool) {
	r := xgb.Rectangle{
		X:      int16(rect.Min.X),
		Y:      int16(rect.Min.Y),
		Width:  uint16(rect.Dx()),
		Height: uint16(rect.Dy()),
	}
	if fill {
		xwin.PolyFillRectangle(win.id, win.gc, []xgb.Rectangle{r})
	}
	xwin.PolyRectangle(win.id, win.gc, []xgb.Rectangle{r})
}

// Lines draws one or more lines in the window.
func (win *window) Lines(pt []image.Point) {
	if len(pt) < 2 {
		return
	}
	pts := make([]xgb.Point, len(pt))
	for i, p := range pt {
		pts[i].X, pts[i].Y = int16(p.X), int16(p.Y)
	}
	xwin.PolyLine(xgb.CoordModeOrigin, win.id, win.gc, pts)
}

// Arc draws an arc on the window.
func (win *window) Arc(rect image.Rectangle, angle1, angle2 float64, fill bool) {
	a := xgb.Arc{
		X:      int16(rect.Min.X),
		Y:      int16(rect.Min.Y),
		Width:  uint16(rect.Dx()),
		Height: uint16(rect.Dy()),
		Angle1: int16(toXAngle(angle1)),
		Angle2: int16(toXAngle(angle2)),
	}
	if fill {
		xwin.PolyFillArc(win.id, win.gc, []xgb.Arc{a})
		return
	}
	xwin.PolyArc(win.id, win.gc, []xgb.Arc{a})
}

func toXAngle(rad float64) float64 {
	return (360 * 64 * rad) / (math.Pi * 2)
}

// Polygon draws a polygon on the window.
func (win *window) Polygon(pt []image.Point, fill bool) {
	if len(pt) < 2 {
		return
	}
	pts := make([]xgb.Point, len(pt))
	for i, p := range pt {
		pts[i].X, pts[i].Y = int16(p.X), int16(p.Y)
	}
	if !fill {
		if !pt[0].Eq(pt[len(pt)-1]) {
			pts = append(pts, xgb.Point{int16(pt[0].X), int16(pt[0].Y)})
		}
		xwin.PolyLine(xgb.CoordModeOrigin, win.id, win.gc, pts)
		return
	}
	xwin.FillPoly(win.id, win.gc, xgb.PolyShapeComplex, xgb.CoordModeOrigin, pts)
}

// Pixel draws a pixel on the window.
func (win *window) Pixel(pt image.Point) {
	xwin.PolyPoint(xgb.CoordModeOrigin, win.id, win.gc, []xgb.Point{xgb.Point{int16(pt.X), int16(pt.Y)}})
}

// SetColor sets the fore color of the window.
func (win *window) SetColor(p sparta.Property, c color.RGBA) {
	if (p != sparta.Background) && (p != sparta.Foreground) {
		return
	}
	s := xwin.DefaultScreen()
	code := getColorCode(c)
	px, ok := pixelMap[code]
	if !ok {
		r, g, b, _ := c.RGBA()
		cl, _ := xwin.AllocColor(s.DefaultColormap, uint16(r), uint16(g), uint16(b))
		px = cl.Pixel
		pixelMap[code] = px
	}
	if p == sparta.Foreground {
		xwin.ChangeGC(win.gc, xgb.GCForeground, []uint32{px})
	} else {
		xwin.ChangeGC(win.gc, xgb.GCBackground, []uint32{px})
	}
}
