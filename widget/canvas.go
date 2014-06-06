// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reservec.
// Distributed under BSD2 license that can be found in LICENSE file.

package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// Arc is an arc for the canvas widget.
type Arc struct {
	Rect           image.Rectangle
	Angle1, Angle2 float64 // in radians
	Fill           bool
}

// Polygon is a polygon for the canvas widget.
type Polygon struct {
	Pt   []image.Point
	Fill bool
}

// Rectangle is a rectangle for the canvas widget.
type Rectangle struct {
	Rect image.Rectangle
	Fill bool
}

// Text is a text for the canvas widget.
type Text struct {
	Pos  image.Point
	Text string
}

var backColor = color.RGBA{R: 255, G: 255, B: 255}
var foreColor = color.RGBA{R: 0, G: 0, B: 0}

// Canvas is a widget in which the client code can draw text, lines,
// rectangles, etc.
type Canvas struct {
	name       string
	win        sparta.Window
	parent     sparta.Widget
	childs     []sparta.Widget
	geometry   image.Rectangle
	fore, back color.RGBA
	border     bool
	data       interface{}

	onDraw, onExpose bool

	closeFn  func(sparta.Widget, interface{}) bool
	commFn   func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	keyFn    func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// NewCanvas creates a new canvas at a given position.
func NewCanvas(parent sparta.Widget, name string, rect image.Rectangle) *Canvas {
	c := &Canvas{
		name:     name,
		parent:   parent,
		geometry: rect,
		back:     backColor,
		fore:     foreColor,
	}
	sparta.NewWindow(c)
	return c
}

// SetWindow is used by the backend to set the backend window of the canvas.
func (c *Canvas) SetWindow(win sparta.Window) {
	c.win = win
}

// Window returns the backend window.
func (c *Canvas) Window() sparta.Window {
	return c.win
}

// RemoveWindow removes the backend window.
func (c *Canvas) RemoveWindow() {
	c.win = nil
}

// Property returns the indicated property of the canvas.
func (c *Canvas) Property(p sparta.Property) interface{} {
	switch p {
	case sparta.Childs:
		return c.childs
	case sparta.Data:
		return c.data
	case sparta.Geometry:
		return c.geometry
	case sparta.Parent:
		return c.parent
	case sparta.Name:
		return c.name
	case sparta.Foreground:
		return c.fore
	case sparta.Background:
		return c.back
	case sparta.Border:
		return c.border
	}
	return nil
}

// SetProperty sets a property of the canvas.
func (c *Canvas) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Childs:
		if v == nil {
			c.childs = nil
			return
		}
		c.childs = append(c.childs, v.(sparta.Widget))
	case sparta.Data:
		c.data = v
	case sparta.Geometry:
		val := v.(image.Rectangle)
		if !c.geometry.Eq(val) {
			c.win.SetProperty(sparta.Geometry, val)
		}
	case sparta.Parent:
		if v == nil {
			c.parent = nil
		}
	case sparta.Name:
		val := v.(string)
		if c.name != val {
			c.name = val
		}
	case sparta.Foreground:
		val := v.(color.RGBA)
		if c.fore != val {
			c.fore = val
			c.win.SetProperty(sparta.Foreground, val)
		}
	case sparta.Background:
		val := v.(color.RGBA)
		if c.back != val {
			c.back = val
			c.win.SetProperty(sparta.Background, val)
		}
	case sparta.Border:
		val := v.(bool)
		if c.border != val {
			c.border = val
		}
	}
}

// Capture sets an event function of the canvas.
func (c *Canvas) Capture(e sparta.EventType, fn func(sparta.Widget, interface{}) bool) {
	switch e {
	case sparta.CloseEv:
		c.closeFn = fn
	case sparta.Configure:
		c.configFn = fn
	case sparta.Command:
		c.commFn = fn
	case sparta.Expose:
		c.exposeFn = fn
	case sparta.KeyEv:
		c.keyFn = fn
	case sparta.Mouse:
		c.mouseFn = fn
	}
}

// OnEvent process a particular event on the canvas.
func (c *Canvas) OnEvent(e interface{}) {
	switch e.(type) {
	case sparta.CloseEvent:
		if c.closeFn != nil {
			c.closeFn(c, e)
		}
		for _, ch := range c.childs {
			ch.OnEvent(e)
		}
	case sparta.ConfigureEvent:
		c.geometry = e.(sparta.ConfigureEvent).Rect
		if c.configFn != nil {
			c.configFn(c, e)
		}
	case sparta.CommandEvent:
		if c.commFn != nil {
			if c.commFn(c, e) {
				return
			}
		}
		c.parent.OnEvent(e)
	case sparta.ExposeEvent:
		if c.exposeFn != nil {
			c.onExpose = true
			c.onDraw = true
			c.exposeFn(c, e)
			c.onExpose = false
			c.onDraw = false
		}
		for _, ch := range c.childs {
			ch.Update()
		}
		if c.border {
			c.win.SetColor(sparta.Foreground, foreColor)
			rect := image.Rect(0, 0, c.geometry.Dx()-1, c.geometry.Dy()-1)
			c.win.Rectangle(rect, false)
		}
	case sparta.KeyEvent:
		if c.keyFn != nil {
			if c.keyFn(c, e) {
				return
			}
		}
		c.parent.OnEvent(e)
	case sparta.MouseEvent:
		if c.mouseFn != nil {
			if c.mouseFn(c, e) {
				return
			}
		}
		ev := e.(sparta.MouseEvent)
		ev.Loc = ev.Loc.Add(c.geometry.Min)
		c.parent.OnEvent(ev)
	}
}

// Update updates the canvas.
func (c *Canvas) Update() {
	c.win.Update()
}

// Focus set the focus on the canvas.
func (c *Canvas) Focus() {
	c.win.Focus()
}

// DrawMode starts the drawing process of the canvas. Before any drawing
// outside of an expose event, this function must be called with the mode
// parameter as true. When drawing functions end, then, the DrawMode function
// must be called with false mode.
func (c *Canvas) DrawMode(mode bool) {
	if c.onExpose {
		return
	}
	if c.onDraw == mode {
		return
	}
	c.win.Draw(mode)
	c.onDraw = mode
}

// Draw draws an object in the canvas, in the given position.
func (c *Canvas) Draw(v interface{}) {
	switch obj := v.(type) {
	case Text:
		c.win.Text(obj.Pos, obj.Text)
	case Rectangle:
		c.win.Rectangle(obj.Rect, obj.Fill)
	case []image.Point:
		c.win.Lines(obj)
	case image.Point:
		c.win.Pixel(obj)
	case Arc:
		c.win.Arc(obj.Rect, obj.Angle1, obj.Angle2, obj.Fill)
	case Polygon:
		c.win.Polygon(obj.Pt, obj.Fill)
	}
}

// SetColor sets the color of the canvas for the next drawing
// operations. If the canvas receive an expose event, the colors will
// be back to defaults.
func (c *Canvas) SetColor(p sparta.Property, cl color.RGBA) {
	if (p != sparta.Background) && (p != sparta.Foreground) {
		return
	}
	c.win.SetColor(p, cl)
}
