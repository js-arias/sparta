// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

package sparta

import "image"

// A CommandEvent is an event sent from goroutines to particular windows.
type CommandEvent struct {
	Source Widget // widget that issue the command
	Value  int    // value identifier of the event
}

// SendEvent sends a command event to an specified window.
var SendEvent = func(dest Widget, comm CommandEvent) {
	panic("undefined Close in the backend")
}

// A KeyEvent is sent for a key press or release.
type KeyEvent struct {
	// The value k represent key k being pressed.
	// The value -k represents key k being released.
	Key Key

	// State represents the keyboard, button state
	State StateKey

	// Loc is the location of the mouse pointer
	Loc image.Point
}

// A MouseEvent is sent for a button press or release or for a mouse movement.
type MouseEvent struct {
	// Button represents the button being pressed or released (negative).
	Button MouseButton

	// State represents the keyboard, button state
	State StateKey

	// Loc is the location of the mouse pointer.
	Loc image.Point
}

// A ConfigureEvent is sent when the window change its size.
type ConfigureEvent struct {
	Rect image.Rectangle
}

// An ExposeEvent is sent when the window has exposed.
type ExposeEvent struct {
	Rect image.Rectangle
}

// A CloseEvent is sent when the window is closed.
type CloseEvent struct{}

// A EventType is a type of event.
type EventType string

// Event types that a window can receive. They are used to define the
// callbacks.
const (
	CloseEv   EventType = "close"     // close event
	Command             = "command"   // command event
	Configure           = "configure" // configure event
	Expose              = "expose"    // expose event
	KeyEv               = "key"       // key events
	Mouse               = "mouse"     // mouse events
)

// keep the window that blocks the input.
var blocker Widget
var blocked = false

// Block block the application from input.
func Block(b Widget) {
	if blocked {
		return
	}
	blocked = true
	blocker = b
}

// IsBlock returns true if the application is blocked from input.
func IsBlock() bool {
	return blocked
}

// IsBlocker returns true if the widget (or its parent) is blocking the
// input.
func IsBlocker(w Widget) bool {
	if w == blocker {
		return true
	}
	for w != nil {
		p := w.Property(Parent)
		if p == nil {
			break
		}
		w = p.(Widget)
		if w == blocker {
			return true
		}
	}
	return false
}

// Unblock opens the application for user input. The widget requesting
// the unblocking must be the same or a child of the widget that request
// the blocking.
func Unblock(req Widget) {
	if !blocked {
		return
	}
	if req == blocker {
		blocked = false
		blocker = nil
		return
	}
	for req != nil {
		p := req.Property(Parent)
		if p == nil {
			break
		}
		req = p.(Widget)
		if req == blocker {
			blocked = false
			blocker = nil
			return
		}
	}
}
