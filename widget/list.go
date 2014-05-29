// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reservec.
// Distributed under BSD2 license that can be found in LICENSE file.

package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// List is a widget that shows a list of strings, and one of the elements
// can be selected with the mouse. When a element of the list is selected, it
// send an event to its target indicating the index of the selected element,
//  it will be a positive number if the selection is made with the left button
// or negative if it was with the right button.
type List struct {
	name       string
	win        sparta.Window
	parent     sparta.Widget
	geometry   image.Rectangle
	fore, back color.RGBA
	data       interface{}

	sel    int
	list   []string
	target sparta.Widget
	scroll *Scroll

	closeFn  func(sparta.Widget, interface{}) bool
	commFn   func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// List particular properties.
const (
	// sets the string list
	ListList sparta.Property = "list"

	// sets the selected element.
	ListSelect = "select"
)

// NewList creates a new list.
func NewList(parent sparta.Widget, name string, rect image.Rectangle) *List {
	l := &List{
		name:     name,
		parent:   parent,
		geometry: rect,
		back:     backColor,
		fore:     foreColor,
		target:   parent,
		sel:      0,
	}
	sparta.NewWindow(l)
	l.scroll = NewScroll(l, "list"+name+"Scroll", 0, 0, Vertical, image.Rect(rect.Dx()-10, 0, rect.Dx(), rect.Dy()))
	return l
}

// SetWindow is used by the backend to sets the backend window of the list.
func (l *List) SetWindow(win sparta.Window) {
	l.win = win
}

// Window returns the backend window.
func (l *List) Window() sparta.Window {
	return l.win
}

// RemoveWindow removes the backend window.
func (l *List) RemoveWindow() {
	l.win = nil
}

// Property returns the indicated property of the list.
func (l *List) Property(p sparta.Property) interface{} {
	switch p {
	case sparta.Childs:
		return []sparta.Widget{l.scroll}
	case sparta.Data:
		return l.data
	case sparta.Geometry:
		return l.geometry
	case sparta.Parent:
		return l.parent
	case sparta.Name:
		return l.name
	case sparta.Foreground:
		return l.fore
	case sparta.Background:
		return l.back
	case sparta.Target:
		return l.target
	case ListList:
		return l.list
	case ListSelect:
		return l.sel
	}
	return nil
}

// SetProperty sets a property of the list.
func (l *List) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Childs:
		if v == nil {
			l.scroll = nil
		}
	case sparta.Data:
		l.data = v
	case sparta.Geometry:
		val := v.(image.Rectangle)
		if !l.geometry.Eq(val) {
			l.win.SetProperty(sparta.Geometry, val)
		}
	case sparta.Parent:
		if v == nil {
			l.parent = nil
		}
	case sparta.Name:
		val := v.(string)
		if l.name != val {
			l.name = val
		}
	case sparta.Foreground:
		val := v.(color.RGBA)
		if l.fore != val {
			l.fore = val
			l.win.SetProperty(sparta.Foreground, val)
		}
	case sparta.Background:
		val := v.(color.RGBA)
		if l.back != val {
			l.back = val
			l.win.SetProperty(sparta.Background, val)
		}
	case sparta.Target:
		val := v.(sparta.Widget)
		if val == nil {
			val = l.parent
		}
		if l.target == val {
			break
		}
		l.target = val
	case ListList:
		val := v.([]string)
		l.list = val
		l.sel = 0
		l.scroll.SetProperty(ScrollSize, 0)
		l.scroll.SetProperty(ScrollSize, len(val))
		l.scroll.SetProperty(ScrollPage, l.geometry.Dy()/sparta.HeightUnit)
		l.Update()
	case ListSelect:
		val := v.(int)
		sel := val
		if sel < 0 {
			sel = -val
		}
		if sel >= len(l.list) {
			break
		}
		l.sel = sel
		pos := l.scroll.Property(ScrollPos).(int)
		if (pos > l.sel) || (l.sel > (pos + l.geometry.Dy()/sparta.HeightUnit)) {
			l.scroll.SetProperty(ScrollPos, l.sel)
		}
		sparta.SendEvent(l.target, sparta.CommandEvent{Source: l, Value: val})
		l.Update()
	}
}

// Capture sets an event function of the list.
func (l *List) Capture(e sparta.EventType, fn func(sparta.Widget, interface{}) bool) {
	switch e {
	case sparta.CloseEv:
		l.closeFn = fn
	case sparta.Configure:
		l.configFn = fn
	case sparta.Command:
		l.commFn = fn
	case sparta.Expose:
		l.exposeFn = fn
	case sparta.KeyEv:
		l.scroll.Capture(e, fn)
	case sparta.Mouse:
		l.mouseFn = fn
	}
}

// OnEvent process a particular event on the list.
func (l *List) OnEvent(e interface{}) {
	switch e.(type) {
	case sparta.CloseEvent:
		if l.closeFn != nil {
			l.closeFn(l, e)
		}
		l.scroll.OnEvent(e)
	case sparta.ConfigureEvent:
		rect := e.(sparta.ConfigureEvent).Rect
		l.geometry = rect
		if l.configFn != nil {
			l.configFn(l, e)
		}
		l.scroll.SetProperty(sparta.Geometry, image.Rect(rect.Dx()-10, 0, rect.Dx(), rect.Dy()))
		l.scroll.SetProperty(ScrollPage, l.geometry.Dy()/sparta.HeightUnit)
	case sparta.CommandEvent:
		if l.commFn != nil {
			if l.commFn(l, e) {
				return
			}
		}
		ev := e.(sparta.CommandEvent)
		if ev.Source == l.scroll {
			l.Update()
			return
		}
		l.parent.OnEvent(e)
	case sparta.ExposeEvent:
		if l.exposeFn != nil {
			l.exposeFn(l, e)
		}
		l.win.SetColor(sparta.Foreground, foreColor)
		if len(l.list) > 0 {
			pos := l.scroll.Property(ScrollPos).(int)
			if pos < 0 {
				pos = 0
			}
			page := l.scroll.Property(ScrollPage).(int)
			for i, s := range l.list[pos:] {
				l.win.Text(image.Pt(2+sparta.WidthUnit, (i*sparta.HeightUnit)+2), s)
				if i > page {
					break
				}
			}
			if (l.sel >= pos) && ((l.sel - pos) < page) {
				y := ((l.sel - pos) * sparta.HeightUnit) + 2
				l.win.Text(image.Pt(2, y), ">")
			}
		}
		rect := image.Rect(0, 0, l.geometry.Dx()-1, l.geometry.Dy()-1)
		l.win.Rectangle(rect, false)
	case sparta.KeyEvent:
		// Only receive this keys if they are returned by the scroll,
		// as the input focus is set on the scroll.
		l.parent.OnEvent(e)
	case sparta.MouseEvent:
		if l.mouseFn != nil {
			if l.mouseFn(l, e) {
				return
			}
		}
		pos := l.scroll.Property(ScrollPos).(int)
		ev := e.(sparta.MouseEvent)
		switch ev.Button {
		case -sparta.MouseWheel:
			l.scroll.SetProperty(ScrollPos, pos+1)
		case sparta.MouseWheel:
			l.scroll.SetProperty(ScrollPos, pos-1)
		case sparta.MouseLeft:
			p := ((ev.Loc.Y - 2) / sparta.HeightUnit) + pos
			l.SetProperty(ListSelect, p)
		case sparta.MouseRight:
			p := ((ev.Loc.Y - 2) / sparta.HeightUnit) + pos
			l.SetProperty(ListSelect, -p)
		}
	}
}

// Update updates the list.
func (l *List) Update() {
	l.win.Update()
}

// Focus set the focus on the list.
func (l *List) Focus() {
	l.scroll.Focus()
}
