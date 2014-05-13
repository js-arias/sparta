// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

package sparta

// Widget is an specific purpouse window.
type Widget interface {
	// SetWindow is used to set the backend window.
	SetWindow(Window)

	// Window returns the backend window.
	Window() Window

	// RemoveWindow removes the backend window.
	RemoveWindow()

	// Property returns the indicated property of the windget.
	Property(Property) interface{}

	// SetProperty sets a property of the winget.
	SetProperty(Property, interface{})

	// Capture sets an event function of the widget. The function
	// receive a Widget and an interface that represents the event,
	// this function returns false if the default operation should
	// be performed.
	Capture(EventType, func(Widget, interface{}) bool)

	// OnEvent is used to send a particular event to the widget.
	OnEvent(interface{})

	// Update updates the widget.
	Update()

	// Focus set the focus on the widget.
	Focus()
}
