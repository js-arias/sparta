// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// Sparta process different events. To connect an event to a function
// the widget method Capture should be used, indicating the captured
// event, and the destination function. The signature of a event function
// is func(w widget, e interface{}) bool, where w is the widget that receive
// the event, e the event structure, and it should return false if further
// default behavior of the widget will follow.
//
// In this example, a window will report all the events that it receive, and
// report them to the standard output. It is based on several similar
// examples, specially from Johnson, E.F. & Reichard, K. (1990) example 5.1
// in "X Window: Applications programming", MIS: Press.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/js-arias/sparta"
	_ "github.com/js-arias/sparta/init"
	"github.com/js-arias/sparta/widget"
)

func main() {
	rand.Seed(time.Now().Unix())

	// The main window.
	m := widget.NewMainWindow("main", "Event Window")

	// Functions to capture the events.
	//
	// Capture close event.This event is produced when attepting to
	// close a rootless widget (such as MainWindow), or when a child
	// is closed.
	m.Capture(sparta.CloseEv, closeFunc)
	// Capture command event. This event is produced when another
	// process, or a widget send a message to another widget.
	m.Capture(sparta.Command, commandFunc)
	// Capture configure event. This event is produced when a widget
	// changes its size.
	m.Capture(sparta.Configure, configureFunc)
	// Capture Expose event. This event is produced when a previously
	// obscured part of a widget is exposed (in some backends, other
	// backends save obscured part of the widget and simply copy them
	// when the widget is exposed. Expose event are also produced when
	// the widget is resized (configure), and when an explicit update
	// method is called.
	m.Capture(sparta.Expose, exposeFunc)
	// Capture Key event. This event is produced when a key in the
	// keyboard is pressed or released.
	m.Capture(sparta.KeyEv, keyFunc)
	// Capture Mouse event. This event is produced when the mouse is
	// moved, a button is pressed or released.
	m.Capture(sparta.Mouse, mouseFunc)

	// this helper function is used to send command events to the window.
	// it will send a random number each second.
	go func() {
		for {
			time.Sleep(10 * time.Second)
			sparta.SendEvent(m, sparta.CommandEvent{Value: rand.Intn(255)})
		}
	}()

	sparta.Run()
}

// CloseFunc process the close event. This event can be stopped if the widget
// is rootless and returns true. If the receiveng widget is a child widget,
// then the closing is unavoidable.
func closeFunc(m sparta.Widget, e interface{}) bool {
	fmt.Fprintf(os.Stdout, "Ev:%s\n", sparta.CloseEv)
	return false
}

// CommandFunc process command event. This event contains the source of the
// event (in this case, is null, as is not sended by any widget), and a value
// that should be interpretated by the receiving widget.
func commandFunc(m sparta.Widget, e interface{}) bool {
	ev := e.(sparta.CommandEvent)
	src := "null"
	if ev.Source != nil {
		// Get the name of the widget
		src = ev.Source.Property(sparta.Name).(string)
	}
	fmt.Fprintf(os.Stdout, "Ev:%s\tSrc:%s\tVal:%d\n", sparta.Command, src, ev.Value)
	return false
}

// ConfigureFunc process configure event. Before any process, the new dimensions
// of the widget where stored in the sparta.Geometry property of the widget.
// The content of the event is a rectangle with the current dimensions and
// position of the widget.
func configureFunc(m sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ConfigureEvent)
	fmt.Fprintf(os.Stdout, "Ev:%s\tMinX:%d\tMinY:%d\tMaxX:%d\tMaxY:%d\n", sparta.Configure, ev.Rect.Min.X, ev.Rect.Min.Y, ev.Rect.Max.X, ev.Rect.Max.Y)
	return false
}

// ExposeFunc process expose event. The content of the event is the invalid
// rectangle that must be restored.
func exposeFunc(m sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ExposeEvent)
	fmt.Fprintf(os.Stdout, "Ev:%s\tMinX:%d\tMinY:%d\tMaxX:%d\tMaxY:%d\n", sparta.Expose, ev.Rect.Min.X, ev.Rect.Min.Y, ev.Rect.Max.X, ev.Rect.Max.Y)
	return false
}

// KeyFunc process the key event. The content of the event is a key code, the
// key code is positive if the key is press, and negative if the key is
// released. If the key is a valid character, it can be used as a rune char.
func keyFunc(m sparta.Widget, e interface{}) bool {
	ev := e.(sparta.KeyEvent)
	press := "press"
	if ev.Key < 0 {
		press = "release"
	}
	r := ""
	if (ev.Key > 0) && ((ev.Key & sparta.KeyNoChar) == 0) {
		r = string([]rune{rune(ev.Key)})
	}
	if len(r) > 0 {
		fmt.Fprintf(os.Stdout, "Ev:%s\t%s\tVal:%s\tSt:%d\tX:%d\tY:%d\n", sparta.KeyEv, press, r, ev.State, ev.Loc.X, ev.Loc.Y)
	} else {
		fmt.Fprintf(os.Stdout, "Ev:%s\t%s\tVal:%d\tSt:%d\tX:%d\tY:%d\n", sparta.KeyEv, press, ev.Key, ev.State, ev.Loc.X, ev.Loc.Y)
	}
	return false
}

// MouseFunc process the mouse event. The content of the event is the button
// action, and the position of the mouse.
func mouseFunc(m sparta.Widget, e interface{}) bool {
	ev := e.(sparta.MouseEvent)
	button := ""
	switch ev.Button {
	case sparta.MouseLeft, -sparta.MouseLeft:
		button = "left"
	case sparta.MouseRight, -sparta.MouseRight:
		button = "right"
	case sparta.Mouse2, -sparta.Mouse2:
		button = "center"
	case sparta.MouseWheel:
		button = "up"
	case -sparta.MouseWheel:
		button = "down"
	}
	press := "move"
	if len(button) > 0 {
		if ev.Button > 0 {
			press = "press"
		} else {
			press = "release"
		}
	}
	fmt.Fprintf(os.Stdout, "Ev:%s\t%s\t%s\tSt:%d\tX:%d\tY:%d\n", sparta.Mouse, press, button, ev.State, ev.Loc.X, ev.Loc.Y)
	return false
}
