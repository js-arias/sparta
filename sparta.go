// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// Package sparta provides a generic (and spartan) widget toolkit for
// window-based guis.
package sparta

// Run executes the program event loop.
var Run = func() {
	panic("undefined Run in the backend")
}

// Close closes the application.
var Close = func() {
	panic("undefined Close in the backend")
}

// Property is a windget property.
type Property string

// Properties of a Widget-Window. The type associated with the property
// is indicated in parentesis.
const (
	// widget title/caption (string)
	Caption Property = "caption"

	// children of the window([]Widget)
	Childs = "childs"

	// data is the particular data associated with the window
	Data = "data"

	// widget geometry (image.Rectangle)
	Geometry = "geometry"

	// parent widget (Widget)
	Parent = "parent"

	// widget id-name (string)
	Name = "name"

	// Background of the window (color.RGBA). This value will be take
	// efect in the next expose event of the widget.
	Background = "background"

	// Foreground of the window (color.RGBA)  This value will be take
	// efect in the next expose event of the widget.
	Foreground = "foreground"

	// Border sets a border line around a widget (bool). If true the
	// widget perimeter will be enmarked. This value will be take
	// efect in the next expose event of the widget.
	Border = "border"

	// Target widget (Widget), used in widgets that sends events to
	// another widget (such a button). If the target is set to nil, then
	// the widget will send the events to its parent.
	Target = "target"
)

// Sparta generic units
var (
	WidthUnit  int
	HeightUnit int
)
