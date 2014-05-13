// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build windows

package win32

import (
	"image"
	"image/color"
	"math"

	"github.com/AllenDang/w32"
	"github.com/js-arias/sparta"
)

type brush struct {
	color w32.COLORREF
	pen   w32.HPEN
	brush w32.HBRUSH
}

func getBrush(c color.RGBA) *brush {
	code := rgb(c)
	b, ok := mapBrush[code]
	if !ok {
		b = &brush{
			color: code,
			pen:   createPen(w32.PS_SOLID, 1, code),
			brush: createSolidBrush(code),
		}
		mapBrush[code] = b
		// BUG: the program will crash if the number
		// of pens is beyond 4990.
	}
	return b
}

// Draw sets the drawing mode of the window.
func (win *window) Draw(mode bool) {
	if win.isPaint {
		return
	}
	if mode {
		if win.dc != 0 {
			return
		}
		win.dc = w32.GetDC(win.id)
		w32.SetBkMode(win.dc, w32.TRANSPARENT)
		w32.SetBkColor(win.dc, win.back.color)
		w32.SelectObject(win.dc, w32.HGDIOBJ(win.fore.brush))
		w32.SelectObject(win.dc, w32.HGDIOBJ(win.fore.pen))
		w32.SelectObject(win.dc, w32.HGDIOBJ(winFont))
		w32.SetTextColor(win.dc, win.fore.color)
		win.curr = win.fore
		return
	}
	if win.dc == 0 {
		return
	}
	w32.ReleaseDC(win.id, win.dc)
	win.dc = 0
}

// Text draws text in the window.
func (win *window) Text(pt image.Point, text string) {
	if len(text) == 0 {
		return
	}
	if win.dc == 0 {
		return
	}
	w32.SetBkMode(win.dc, w32.OPAQUE)
	textOut(win.dc, pt.X, pt.Y, text)
	w32.SetBkMode(win.dc, w32.TRANSPARENT)
}

// Rectangle draws a rectangle in the window.
func (win *window) Rectangle(rect image.Rectangle, fill bool) {
	if win.dc == 0 {
		return
	}
	if !fill {
		pts := []w32.POINT{
			{X: int32(rect.Min.X), Y: int32(rect.Min.Y)},
			{X: int32(rect.Max.X), Y: int32(rect.Min.Y)},
			{X: int32(rect.Max.X), Y: int32(rect.Max.Y)},
			{X: int32(rect.Min.X), Y: int32(rect.Max.Y)},
			{X: int32(rect.Min.X), Y: int32(rect.Min.Y)},
		}
		polyLine(win.dc, pts)
		return
	}
	w32.Rectangle(win.dc, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

// Lines draws one or more lines in the window.
func (win *window) Lines(pt []image.Point) {
	if win.dc == 0 {
		return
	}
	if len(pt) < 2 {
		return
	}
	pts := make([]w32.POINT, len(pt))
	for i, p := range pt {
		pts[i].X, pts[i].Y = int32(p.X), int32(p.Y)
	}
	polyLine(win.dc, pts)
}

// Arc draws an arc on the window.
func (win *window) Arc(rect image.Rectangle, angle1, angle2 float64, fill bool) {
	if win.dc == 0 {
		return
	}
	xStart := int(math.Cos(angle1)*(float64(rect.Dx())/2)) + rect.Min.X
	xEnd := int(math.Cos(angle2)*(float64(rect.Dx())/2)) + rect.Min.X
	yStart := int(math.Sin(angle1)*(float64(rect.Dy())/2)) + rect.Min.Y
	yEnd := int(math.Sin(angle2)*(float64(rect.Dy())/2)) + rect.Min.Y
	if fill {
		pie(win.dc, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, int(xStart), int(yStart), int(xEnd), int(yEnd))
		return
	}
	arc(win.dc, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, int(xStart), int(yStart), int(xEnd), int(yEnd))
}

// Polygon draws a polygon on the window.
func (win *window) Polygon(pt []image.Point, fill bool) {
	if len(pt) < 2 {
		return
	}
	pts := make([]w32.POINT, len(pt))
	for i, p := range pt {
		pts[i].X, pts[i].Y = int32(p.X), int32(p.Y)
	}
	if !fill {
		if !pt[0].Eq(pt[len(pt)-1]) {
			pts = append(pts, w32.POINT{pts[0].X, pts[0].Y})
		}
		polyLine(win.dc, pts)
		return
	}
	polygon(win.dc, pts)
}

// Pixel draws a pixel on the window.
func (win *window) Pixel(pt image.Point) {
	if win.dc == 0 {
		return
	}
	setPixel(win.dc, pt.X, pt.Y, win.curr.color)
}

// SetColor sets the drawing color of the window.
func (win *window) SetColor(p sparta.Property, c color.RGBA) {
	if (p != sparta.Background) && (p != sparta.Foreground) {
		return
	}
	if win.dc == 0 {
		return
	}
	b := getBrush(c)
	if p == sparta.Foreground {
		w32.SelectObject(win.dc, w32.HGDIOBJ(b.brush))
		w32.SelectObject(win.dc, w32.HGDIOBJ(b.pen))
		w32.SetTextColor(win.dc, b.color)
		win.curr = b
	} else {
		w32.SetBkColor(win.dc, b.color)
	}
}
