// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

package sparta

import (
	"image"
	"image/color"
)

// Drawable is a screen space to be draw by the backend.
type Drawable interface {
	// Draw sets a drawable in a ready to draw mode (true) or ends the
	// drawing mode (false).
	Draw(bool)

	// SetColor sets a temporal drawing color of the drawable. This
	// values will be valid until a new Expose event will be produced.
	// If you want to change the values of the window, use "SetProperty".
	SetColor(Property, color.RGBA)

	// Text draws a string in the indicated point of the drawable.
	Text(image.Point, string)

	// Rectangle draws a rectangle in the drawable.
	Rectangle(image.Rectangle, bool)

	// Lines draws one or more lines in the drawable.
	Lines([]image.Point)

	// Arc draws an arc in the drawable.
	Arc(image.Rectangle, float64, float64, bool)

	// Polygon draws a polygon in the drawable.
	Polygon([]image.Point, bool)

	// Pixel draws a pixel in the drawable.
	Pixel(image.Point)
}

// Window is a backend window that receives events. Most code should
// use Widget type.
type Window interface {
	Drawable

	// Close closes the window.
	Close()

	// SetProperty sets a window property.
	SetProperty(Property, interface{})

	// Update updates the window content.
	Update()

	// Focus set the focus on the window.
	Focus()
}

// NewWindow assigns a new window to a widget.
var NewWindow = func(w Widget) {
	panic("undefined NewWindow in the backend")
}
