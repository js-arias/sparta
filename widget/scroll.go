// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reservec.
// Distributed under BSD2 license that can be found in LICENSE file.

package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// ScrollType is the kind of the scroll
type ScrollType bool

const (
	// Vertical scroll
	Vertical ScrollType = false

	// Horizontal scroll
	Horizontal = true
)

// Scroll is a widget that shows a position inside a document. When a scroll
// is moved, it sends an event to its target indicating the new
// position.
//
// If you are using an scroll, and want to move the content of the target
// client, the change the scroll position property, and process the movement
// when the target client receive the position event from the scroll.
type Scroll struct {
	name       string
	win        sparta.Window
	parent     sparta.Widget
	geometry   image.Rectangle
	fore, back color.RGBA
	data       interface{}

	pos, size, page int
	typ             ScrollType
	target          sparta.Widget

	closeFn  func(sparta.Widget, interface{}) bool
	commFn   func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	keyFn    func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// Scroll particular properties.
const (
	// sets the size of the scroll page (int)
	ScrollPage sparta.Property = "page"

	// sets the position of the scroll (int)
	ScrollPos = "pos"

	// sets the size of the scroll (int)
	ScrollSize = "size"
)

// NewScroll creates a new scroll of the given type.
func NewScroll(parent sparta.Widget, name string, size, page int, typ ScrollType, rect image.Rectangle) *Scroll {
	s := &Scroll{
		name:     name,
		parent:   parent,
		geometry: rect,
		back:     backColor,
		fore:     foreColor,
		size:     size,
		page:     page,
		target:   parent,
		typ:      typ,
	}
	sparta.NewWindow(s)
	return s
}

// SetWindow is used by the backend to sets the backend window of the scroll.
func (s *Scroll) SetWindow(win sparta.Window) {
	s.win = win
}

// Window returns the backend window.
func (s *Scroll) Window() sparta.Window {
	return s.win
}

// RemoveWindow removes the backend window.
func (s *Scroll) RemoveWindow() {
	s.win = nil
}

// Property returns the indicated property of the scroll.
func (s *Scroll) Property(p sparta.Property) interface{} {
	switch p {
	case sparta.Data:
		return s.data
	case sparta.Geometry:
		return s.geometry
	case sparta.Parent:
		return s.parent
	case sparta.Name:
		return s.name
	case sparta.Foreground:
		return s.fore
	case sparta.Background:
		return s.back
	case sparta.Target:
		return s.target
	case ScrollPage:
		return s.page
	case ScrollPos:
		return s.pos
	case ScrollSize:
		return s.size
	}
	return nil
}

// SetProperty sets a property of the scroll.
func (s *Scroll) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Data:
		s.data = v
	case sparta.Geometry:
		val := v.(image.Rectangle)
		if !s.geometry.Eq(val) {
			s.win.SetProperty(sparta.Geometry, val)
		}
	case sparta.Parent:
		if v == nil {
			s.parent = nil
		}
	case sparta.Name:
		val := v.(string)
		if s.name != val {
			s.name = val
		}
	case sparta.Foreground:
		val := v.(color.RGBA)
		if s.fore != val {
			s.fore = val
			s.win.SetProperty(sparta.Foreground, val)
		}
	case sparta.Background:
		val := v.(color.RGBA)
		if s.back != val {
			s.back = val
			s.win.SetProperty(sparta.Background, val)
		}
	case sparta.Target:
		val := v.(sparta.Widget)
		if val == nil {
			val = s.parent
		}
		if s.target == val {
			break
		}
		s.target = val
	case ScrollPage:
		val := v.(int)
		if val < 0 {
			break
		}
		if s.page == val {
			break
		}
		s.page = val
		if ((s.size - s.page) > 0) && (s.pos > (s.size - s.page)) {
			s.pos = s.size - s.page
			sparta.SendEvent(s.target, sparta.CommandEvent{Source: s, Value: s.pos})
		}
		s.Update()
	case ScrollPos:
		val := v.(int)
		if val < 0 {
			val = 0
		} else if val > (s.size - s.page) {
			val = s.size - s.page
		}
		if s.pos == val {
			break
		}
		s.pos = val
		sparta.SendEvent(s.target, sparta.CommandEvent{Source: s, Value: s.pos})
		s.Update()
	case ScrollSize:
		val := v.(int)
		if val < 0 {
			break
		}
		if s.size == val {
			break
		}
		s.size = val
		if s.size == 0 {
			s.page = 0
			s.pos = 0
		} else if s.pos > (s.size - s.page) {
			s.pos = s.size - s.page
			sparta.SendEvent(s.target, sparta.CommandEvent{Source: s, Value: s.pos})
		}
		s.Update()
	}
}

// Capture sets an event function of the scroll.
func (s *Scroll) Capture(e sparta.EventType, fn func(sparta.Widget, interface{}) bool) {
	switch e {
	case sparta.CloseEv:
		s.closeFn = fn
	case sparta.Configure:
		s.configFn = fn
	case sparta.Command:
		s.commFn = fn
	case sparta.Expose:
		s.exposeFn = fn
	case sparta.KeyEv:
		s.keyFn = fn
	case sparta.Mouse:
		s.mouseFn = fn
	}
}

// OnEvent process a particular event on the scroll.
func (s *Scroll) OnEvent(e interface{}) {
	switch e.(type) {
	case sparta.CloseEvent:
		if s.closeFn != nil {
			s.closeFn(s, e)
		}
	case sparta.ConfigureEvent:
		s.geometry = e.(sparta.ConfigureEvent).Rect
		if s.configFn != nil {
			s.configFn(s, e)
		}
	case sparta.CommandEvent:
		if s.commFn != nil {
			if s.commFn(s, e) {
				return
			}
		}
		s.parent.OnEvent(e)
	case sparta.ExposeEvent:
		if s.exposeFn != nil {
			s.exposeFn(s, e)
		}
		s.win.SetColor(sparta.Foreground, foreColor)
		rect := image.Rect(0, 0, s.geometry.Dx()-1, s.geometry.Dy()-1)
		s.win.Rectangle(rect, false)
		if s.size > 0 {
			if s.typ == Vertical {
				rect.Min.Y = (s.geometry.Dy() * s.pos) / s.size
				rect.Max.Y = rect.Min.Y + ((s.geometry.Dy() * s.page) / s.size)
			} else {
				rect.Min.X = (s.geometry.Dx() * s.pos) / s.size
				rect.Max.X = rect.Min.X + ((s.geometry.Dx() * s.page) / s.size)
			}
			s.win.Rectangle(rect, true)
		}
	case sparta.KeyEvent:
		if s.keyFn != nil {
			if s.keyFn(s, e) {
				return
			}
		}
		ev := e.(sparta.KeyEvent)
		switch ev.Key {
		case sparta.KeyDown:
			if s.typ == Vertical {
				s.SetProperty(ScrollPos, s.pos+1)
				return
			}
		case sparta.KeyUp:
			if s.typ == Vertical {
				s.SetProperty(ScrollPos, s.pos-1)
				return
			}
		case sparta.KeyLeft:
			if s.typ != Vertical {
				s.SetProperty(ScrollPos, s.pos-1)
				return
			}
		case sparta.KeyRight:
			if s.typ != Vertical {
				s.SetProperty(ScrollPos, s.pos+1)
				return
			}
		case sparta.KeyPageUp:
			if s.typ == Vertical {
				s.SetProperty(ScrollPos, s.pos-s.page)
				return
			}
		case sparta.KeyPageDown:
			if s.typ == Vertical {
				s.SetProperty(ScrollPos, s.pos+s.page)
				return
			}
		case sparta.KeyHome:
			s.SetProperty(ScrollPos, 0)
			return
		case sparta.KeyEnd:
			s.SetProperty(ScrollPos, s.size)
			return
		}
		s.parent.OnEvent(e)
	case sparta.MouseEvent:
		if s.mouseFn != nil {
			if s.mouseFn(s, e) {
				return
			}
		}
		ev := e.(sparta.MouseEvent)
		switch ev.Button {
		case sparta.MouseWheel:
			s.SetProperty(ScrollPos, s.pos-1)
		case -sparta.MouseWheel:
			s.SetProperty(ScrollPos, s.pos+1)
		case sparta.MouseLeft:
			if s.typ == Vertical {
				p := (ev.Loc.Y * s.size) / s.geometry.Dy()
				s.SetProperty(ScrollPos, p)
			} else {
				p := (ev.Loc.X * s.size) / s.geometry.Dx()
				s.SetProperty(ScrollPos, p)
			}
		case sparta.MouseRight:
			if s.typ == Vertical {
				p := (ev.Loc.Y * s.size) / s.geometry.Dy()
				s.SetProperty(ScrollPos, p-s.page)
			} else {
				p := (ev.Loc.X * s.size) / s.geometry.Dx()
				s.SetProperty(ScrollPos, p-s.page)
			}
		}
	}
}

// Update updates the scroll.
func (s *Scroll) Update() {
	s.win.Update()
}

// Focus set the focus on the scroll.
func (s *Scroll) Focus() {
	s.win.Focus()
}
