// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reservec.
// Distributed under BSD2 license that can be found in LICENSE file.

package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// Button is a widget that shows a text, and can be "pushed" with the mouse.
// When a mouse button is pressed over a button widget, it will sends an
// arbitrary value (that can be set with the propery ButtonValue) to the
// target widget.
type Button struct {
	name       string
	win        sparta.Window
	parent     sparta.Widget
	geometry   image.Rectangle
	fore, back color.RGBA
	data       interface{}

	caption string
	target  sparta.Widget
	value   int

	closeFn  func(sparta.Widget, interface{}) bool
	commFn   func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	keyFn    func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// Button particular properties.
const (
	// sets the value (int) that the button will send to the target if
	// pressed.
	ButtonValue sparta.Property = "value"
)

// Button creates a new button.
func NewButton(parent sparta.Widget, name, caption string, rect image.Rectangle) *Button {
	b := &Button{
		name:     name,
		parent:   parent,
		geometry: rect,
		back:     backColor,
		fore:     foreColor,
		caption:  caption,
		target:   parent,
	}
	sparta.NewWindow(b)
	return b
}

// SetWindow is used by the backend to sets the backend window of the button.
func (b *Button) SetWindow(win sparta.Window) {
	b.win = win
}

// Window returns the backend window.
func (b *Button) Window() sparta.Window {
	return b.win
}

// RemoveWindow removes the backend window.
func (b *Button) RemoveWindow() {
	b.win = nil
}

// Property returns the indicated property of the button.
func (b *Button) Property(p sparta.Property) interface{} {
	switch p {
	case sparta.Caption:
		return b.caption
	case sparta.Data:
		return b.data
	case sparta.Geometry:
		return b.geometry
	case sparta.Parent:
		return b.parent
	case sparta.Name:
		return b.name
	case sparta.Foreground:
		return b.fore
	case sparta.Background:
		return b.back
	case sparta.Target:
		return b.target
	case ButtonValue:
		return b.value
	}
	return nil
}

// SetProperty sets a property of the button.
func (b *Button) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Caption:
		val := v.(string)
		if b.caption != val {
			b.caption = val
			b.Update()
		}
	case sparta.Data:
		b.data = v
	case sparta.Geometry:
		val := v.(image.Rectangle)
		if !b.geometry.Eq(val) {
			b.win.SetProperty(sparta.Geometry, val)
		}
	case sparta.Parent:
		if v == nil {
			b.parent = nil
		}
	case sparta.Name:
		val := v.(string)
		if b.name != val {
			b.name = val
		}
	case sparta.Foreground:
		val := v.(color.RGBA)
		if b.fore != val {
			b.fore = val
			b.win.SetProperty(sparta.Foreground, val)
		}
	case sparta.Background:
		val := v.(color.RGBA)
		if b.back != val {
			b.back = val
			b.win.SetProperty(sparta.Background, val)
		}
	case sparta.Target:
		val := v.(sparta.Widget)
		if val == nil {
			val = b.parent
		}
		if b.target == val {
			break
		}
		b.target = val
	case ButtonValue:
		val := v.(int)
		if b.value != val {
			b.value = val
		}
	}
}

// Capture sets an event function of the button.
func (b *Button) Capture(e sparta.EventType, fn func(sparta.Widget, interface{}) bool) {
	switch e {
	case sparta.CloseEv:
		b.closeFn = fn
	case sparta.Configure:
		b.configFn = fn
	case sparta.Command:
		b.commFn = fn
	case sparta.Expose:
		b.exposeFn = fn
	case sparta.KeyEv:
		b.keyFn = fn
	case sparta.Mouse:
		b.mouseFn = fn
	}
}

// OnEvent process a particularevent on the button.
func (b *Button) OnEvent(e interface{}) {
	switch e.(type) {
	case sparta.CloseEvent:
		if b.closeFn != nil {
			b.closeFn(b, e)
		}
	case sparta.ConfigureEvent:
		b.geometry = e.(sparta.ConfigureEvent).Rect
		if b.configFn != nil {
			b.configFn(b, e)
		}
	case sparta.CommandEvent:
		if b.commFn != nil {
			if b.commFn(b, e) {
				return
			}
		}
		b.parent.OnEvent(e)
	case sparta.ExposeEvent:
		if b.exposeFn != nil {
			b.exposeFn(b, e)
		}
		b.win.SetColor(sparta.Foreground, foreColor)
		if len(b.caption) > 0 {
			x := (b.geometry.Dx() - (len(b.caption) * sparta.WidthUnit)) / 2
			y := (b.geometry.Dy() - sparta.HeightUnit) / 2
			b.win.Text(image.Pt(x, y), b.caption)
		}
		rect := image.Rect(0, 0, b.geometry.Dx()-1, b.geometry.Dy()-1)
		b.win.Rectangle(rect, false)
	case sparta.KeyEvent:
		if sparta.IsBlock() {
			if !sparta.IsBlocker(b) {
				return
			}
		}
		if b.keyFn != nil {
			if b.keyFn(b, e) {
				return
			}
		}
	case sparta.MouseEvent:
		if sparta.IsBlock() {
			if !sparta.IsBlocker(b) {
				return
			}
		}
		if b.mouseFn != nil {
			if b.mouseFn(b, e) {
				return
			}
		}
		ev := e.(sparta.MouseEvent)
		if ev.Button == sparta.MouseLeft {
			sparta.SendEvent(b.target, sparta.CommandEvent{Source: b, Value: b.value})
		}
	}
}

// Update updates the button.
func (b *Button) Update() {
	b.win.Update()
}

// Focus set the focus on the button.
func (b *Button) Focus() {
	b.win.Focus()
}
