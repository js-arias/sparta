// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reservec.
// Distributed under BSD2 license that can be found in LICENSE file.

package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// ListData is a list used by a list control.
type ListData interface {
	// Len returns the length of the list.
	Len() int

	// Item returns the name of the i-th element of the list.
	Item(int) string

	// IsSel returns true if the i-th element is selected.
	IsSel(int) bool
}

// List is a widget that shows a list of strings, and one or more elements
// can be selected with the mouse.
//
// When a element of the list is selected, it send a comand event to its target
// widget indicating the index of the selected element. It will be a positive
// number if the selection is made with the left button or negative (starting
// at -1), if it was with the right button.
//
// It is up to client code to manage multiple or single selection.
type List struct {
	name       string
	win        sparta.Window
	parent     sparta.Widget
	geometry   image.Rectangle
	fore, back color.RGBA
	data       interface{}

	list   ListData
	target sparta.Widget
	scroll *Scroll

	closeFn  func(sparta.Widget, interface{}) bool
	commFn   func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	keyFn    func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// List particular properties.
const (
	// sets the string list
	ListList sparta.Property = "list"
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
		if v == nil {
			l.list = nil
			l.scroll.SetProperty(ScrollSize, 0)
		} else {
			val := v.(ListData)
			l.list = val
			l.scroll.SetProperty(ScrollSize, 0)
			l.scroll.SetProperty(ScrollSize, val.Len())
		}
		l.scroll.SetProperty(ScrollPage, l.geometry.Dy()/sparta.HeightUnit)
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
		l.keyFn = fn
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
		if (l.list != nil) && (l.list.Len() > 0) {
			pos := l.scroll.Property(ScrollPos).(int)
			if pos < 0 {
				pos = 0
			}
			page := l.scroll.Property(ScrollPage).(int)
			for i := 0; i <= page; i++ {
				j := i + pos
				if j >= l.list.Len() {
					break
				}
				if l.list.IsSel(j) {
					y := (i * sparta.HeightUnit) + 2
					l.win.Text(image.Pt(2, y), ">")
				}
				l.win.Text(image.Pt(2+sparta.WidthUnit, (i*sparta.HeightUnit)+2), l.list.Item(j))
			}
		}
		rect := image.Rect(0, 0, l.geometry.Dx()-1, l.geometry.Dy()-1)
		l.win.Rectangle(rect, false)
	case sparta.KeyEvent:
		if s.keyFn != nil {
			if s.keyFn(s, e) {
				return
			}
		}
		pos := l.scroll.Property(ScrollPos).(int)
		page := l.scroll.Property(ScrollPage).(int)
		ev := e.(sparta.KeyEvent)
		switch ev.Key {
		case sparta.KeyDown:
			l.scroll.SetProperty(ScrollPos, pos+1)
		case sparta.KeyUp:
			l.scroll.SetProperty(ScrollPos, pos-1)
		case sparta.KeyPageUp:
			l.scroll.SetProperty(ScrollPos, pos-page)
		case sparta.KeyPageDown:
			l.scroll.SetProperty(ScrollPos, pos+page)
		case sparta.KeyHome:
			l.scroll.SetProperty(ScrollPos, 0)
		case sparta.KeyEnd:
			s.SetProperty(ScrollPos, l.list.Len())
		default:
			l.parent.OnEvent(e)
		}
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
			sparta.SendEvent(l.target, sparta.CommandEvent{Source: l, Value: p})
		case sparta.MouseRight:
			p := ((ev.Loc.Y - 2) / sparta.HeightUnit) + pos
			sparta.SendEvent(l.target, sparta.CommandEvent{Source: l, Value: -(p + 1)})
		}
	}
}

// Update updates the list.
func (l *List) Update() {
	l.win.Update()
}

// Focus set the focus on the list.
func (l *List) Focus() {
	l.win.Focus()
}
