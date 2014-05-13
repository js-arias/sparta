// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// Package widget provides a concrete widgets using the sparta package.
package widget

import (
	"image"
	"image/color"

	"github.com/js-arias/sparta"
)

// MainWindow is a main window.
type MainWindow struct {
	name       string
	win        sparta.Window
	geometry   image.Rectangle
	childs     []sparta.Widget
	fore, back color.RGBA
	data       interface{}

	title string

	commFn   func(sparta.Widget, interface{}) bool
	closeFn  func(sparta.Widget, interface{}) bool
	configFn func(sparta.Widget, interface{}) bool
	exposeFn func(sparta.Widget, interface{}) bool
	keyFn    func(sparta.Widget, interface{}) bool
	mouseFn  func(sparta.Widget, interface{}) bool
}

// NewMainWindow creates a new main window.
func NewMainWindow(name, title string) *MainWindow {
	w := &MainWindow{
		name:     name,
		geometry: image.Rect(0, 0, 80*sparta.WidthUnit, 20*sparta.HeightUnit),
		back:     backColor,
		fore:     foreColor,
		title:    title,
	}
	sparta.NewWindow(w)
	w.win.SetProperty(sparta.Caption, w.title)
	return w
}

// SetWindow sets the backend window of the window.
func (w *MainWindow) SetWindow(win sparta.Window) {
	w.win = win
}

// Window returns the backend window.
func (w *MainWindow) Window() sparta.Window {
	return w.win
}

// RemoveWindow removes the backend window.
func (w *MainWindow) RemoveWindow() {
	w.win = nil
}

// Property returns the indicated property of the main window.
func (w *MainWindow) Property(p sparta.Property) interface{} {
	switch p {
	case sparta.Caption:
		return w.title
	case sparta.Childs:
		return w.childs
	case sparta.Data:
		return w.data
	case sparta.Geometry:
		return w.geometry
	case sparta.Name:
		return w.name
	case sparta.Foreground:
		return w.fore
	case sparta.Background:
		return w.back
	}
	return nil
}

// SetProperty sets a property of the main window.
func (w *MainWindow) SetProperty(p sparta.Property, v interface{}) {
	switch p {
	case sparta.Caption:
		val := v.(string)
		if w.title != val {
			w.title = val
			w.win.SetProperty(sparta.Caption, w.title)
		}
	case sparta.Childs:
		if v == nil {
			w.childs = nil
			return
		}
		w.childs = append(w.childs, v.(sparta.Widget))
	case sparta.Data:
		w.data = v
	case sparta.Geometry:
		val := v.(image.Rectangle)
		if !w.geometry.Eq(val) {
			w.win.SetProperty(sparta.Geometry, val)
		}
	case sparta.Name:
		val := v.(string)
		if w.name != val {
			w.name = val
		}
	case sparta.Foreground:
		val := v.(color.RGBA)
		if w.fore != val {
			w.fore = val
			w.win.SetProperty(sparta.Foreground, val)
		}
	case sparta.Background:
		val := v.(color.RGBA)
		if w.back != val {
			w.back = val
			w.win.SetProperty(sparta.Background, val)
		}
	}
}

// Capture sets an event function of the main window.
func (w *MainWindow) Capture(e sparta.EventType, fn func(sparta.Widget, interface{}) bool) {
	switch e {
	case sparta.CloseEv:
		w.closeFn = fn
	case sparta.Command:
		w.commFn = fn
	case sparta.Configure:
		w.configFn = fn
	case sparta.Expose:
		w.exposeFn = fn
	case sparta.KeyEv:
		w.keyFn = fn
	case sparta.Mouse:
		w.mouseFn = fn
	}
}

// OnEvent process a particular event on the main window.
func (w *MainWindow) OnEvent(e interface{}) {
	switch e.(type) {
	case sparta.CloseEvent:
		if sparta.IsBlock() {
			if !sparta.IsBlocker(w) {
				return
			}
		}
		if w.closeFn != nil {
			if w.closeFn(w, e) {
				return
			}
		}
		for _, c := range w.childs {
			c.OnEvent(e)
		}
		w.win.Close()
	case sparta.CommandEvent:
		if w.commFn != nil {
			w.commFn(w, e)
		}
	case sparta.ConfigureEvent:
		w.geometry = e.(sparta.ConfigureEvent).Rect
		if w.configFn != nil {
			w.configFn(w, e)
		}
	case sparta.ExposeEvent:
		if w.exposeFn != nil {
			w.exposeFn(w, e)
		}
	case sparta.KeyEvent:
		if w.keyFn != nil {
			w.keyFn(w, e)
		}
	case sparta.MouseEvent:
		if w.mouseFn != nil {
			w.mouseFn(w, e)
		}
	}
}

// Update updates the main window.
func (w *MainWindow) Update() {
	w.win.Update()
}

// Focus set the focus on the main window.
func (w *MainWindow) Focus() {
	w.win.Focus()
}
