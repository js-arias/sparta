// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build darwin freebsd linux netbsd openbsd

package x11

import (
	"image"
	"log"
	"os"
	"unicode"

	"github.com/js-arias/sparta"
	"github.com/js-arias/xgb"
)

func init() {
	sparta.Run = run
	sparta.Close = closeApp
	sparta.SendEvent = sendEvent
}

var endChan = make(chan struct{})

func put32(buf []byte, v uint32) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

// SendEvent sends an event to the window.
func sendEvent(dest sparta.Widget, comm sparta.CommandEvent) {
	dwin := dest.Window().(*window)
	var id xgb.Id
	if comm.Source != nil {
		id = comm.Source.Window().(*window).id
	}
	event := make([]byte, 32)
	event[0] = xgb.ClientMessage // id of the event
	event[1] = 32                // format
	put32(event[4:], uint32(dwin.id))
	put32(event[8:], uint32(atomMsg))     // message type (client message)
	put32(event[12:], uint32(id))         // sender of the event
	put32(event[16:], uint32(comm.Value)) // value of the event
	xwin.SendEvent(false, dwin.id, 0, event)
}

const allEventMask = xgb.EventMaskKeyPress | xgb.EventMaskKeyRelease |
	xgb.EventMaskButtonPress | xgb.EventMaskButtonRelease |
	xgb.EventMaskPointerMotion | xgb.EventMaskButtonMotion |
	xgb.EventMaskExposure | xgb.EventMaskStructureNotify

// Run runs the x11 event loop.
func run() {
	evChan := make(chan xgb.Event)
	go func() {
		for {
			reply, err := xwin.WaitForEvent()
			if err != nil {
				log.Printf("x11: error: %v\n", err)
				os.Exit(1)
			}
			evChan <- reply
		}
	}()
	for {
		select {
		case e := <-evChan:
			xEvent(e)
		case <-endChan:
			return
		}
	}
}

// CloseApp closes the application.
func closeApp() {
	close(endChan)
}

// XEvent proccess an x11 event.
func xEvent(e xgb.Event) {
	switch event := e.(type) {
	case xgb.ButtonPressEvent:
		w, ok := widgetTable[event.Event]
		if !ok {
			break
		}
		ev := sparta.MouseEvent{
			Button: getButton(event.Detail),
			State:  sparta.StateKey(event.State),
			Loc:    image.Pt(int(event.EventX), int(event.EventY)),
		}
		w.OnEvent(ev)
		w.Focus()
	case xgb.ButtonReleaseEvent:
		if (event.Detail == 4) || (event.Detail == 5) {
			break
		}
		w, ok := widgetTable[event.Event]
		if !ok {
			break
		}
		ev := sparta.MouseEvent{
			Button: -getButton(event.Detail),
			State:  sparta.StateKey(event.State),
			Loc:    image.Pt(int(event.EventX), int(event.EventY)),
		}
		w.OnEvent(ev)
	case xgb.ClientMessageEvent:
		w, ok := widgetTable[event.Window]
		if !ok {
			break
		}
		switch event.Type {
		case atomMsg:
			src := xgb.Id(event.Data.Data32[0])
			val := int(event.Data.Data32[1])
			sw, ok := widgetTable[src]
			if !ok {
				sw = nil
			}
			w.OnEvent(sparta.CommandEvent{Source: sw, Value: val})
			break
		case wmProtocols:
			if w.Property(sparta.Parent) != nil {
				break
			}
			if event.Type != wmProtocols {
				break
			}
			if xgb.Id(event.Data.Data32[0]) != atomDel {
				break
			}
			w.OnEvent(sparta.CloseEvent{})
		}
	case xgb.DestroyNotifyEvent:
		w, ok := widgetTable[event.Window]
		if !ok {
			break
		}
		if w.Property(sparta.Parent) != nil {
			break
		}
		w.OnEvent(sparta.CloseEvent{})
	case xgb.ConfigureNotifyEvent:
		w, ok := widgetTable[event.Window]
		if !ok {
			break
		}
		rect := w.Property(sparta.Geometry).(image.Rectangle)
		if (rect.Dx() == int(event.Width)) && (rect.Dy() == int(event.Height)) {
			break
		}
		ev := sparta.ConfigureEvent{image.Rect(int(event.X), int(event.Y), int(event.X)+int(event.Width), int(event.Y)+int(event.Height))}
		w.OnEvent(ev)
		xwin.ClearArea(true, event.Window, 0, 0, event.Width, event.Height)
	case xgb.ExposeEvent:
		// only proccess the last expose event
		if event.Count != 0 {
			break
		}
		w, ok := widgetTable[event.Window]
		if !ok {
			break
		}
		win := w.Window().(*window)
		xwin.ChangeGC(win.gc, xgb.GCForeground, []uint32{win.back})
		r := xgb.Rectangle{
			X:      int16(event.X),
			Y:      int16(event.Y),
			Width:  event.Width,
			Height: event.Height,
		}
		xwin.PolyFillRectangle(win.id, win.gc, []xgb.Rectangle{r})
		xwin.ChangeGC(win.gc, xgb.GCForeground, []uint32{win.fore})
		win.isExpose = true
		ev := sparta.ExposeEvent{image.Rect(int(event.X), int(event.Y), int(event.X+event.Width), int(event.Y+event.Height))}
		w.OnEvent(ev)
		win.isExpose = false
	case xgb.KeyPressEvent:
		w, ok := widgetTable[event.Event]
		if !ok {
			break
		}
		ev := sparta.KeyEvent{
			Key:   sparta.Key(getKeyValue(int(event.Detail), int(event.State))),
			State: sparta.StateKey(event.State),
			Loc:   image.Pt(int(event.EventX), int(event.EventY)),
		}
		w.OnEvent(ev)
	case xgb.KeyReleaseEvent:
		w, ok := widgetTable[event.Event]
		if !ok {
			break
		}
		ev := sparta.KeyEvent{
			Key:   -sparta.Key(keysyms[int(event.Detail)][0]),
			State: sparta.StateKey(event.State),
			Loc:   image.Pt(int(event.EventX), int(event.EventY)),
		}
		if (ev.Key - 1) == sparta.KeyShift {
			ev.Key = sparta.KeyShift
		}
		if (ev.Key - 1) == sparta.KeyControl {
			ev.Key = sparta.KeyControl
		}
		w.OnEvent(ev)
	case xgb.MappingNotifyEvent:
		setKeyboard()
	case xgb.MotionNotifyEvent:
		w, ok := widgetTable[event.Event]
		if !ok {
			break
		}
		ev := sparta.MouseEvent{
			Button: getButton(event.Detail),
			State:  sparta.StateKey(event.State),
			Loc:    image.Pt(int(event.EventX), int(event.EventY)),
		}
		w.OnEvent(ev)
	}
}

func getButton(button byte) sparta.MouseButton {
	switch button {
	case 1:
		return sparta.MouseLeft
	case 2:
		return sparta.Mouse2
	case 3:
		return sparta.MouseRight
	case 4:
		return sparta.MouseWheel
	case 5:
		return -sparta.MouseWheel
	}
	return 0
}

// GetKeyValue gets keyboard encoding as document at:
// http://tronche.com/gui/x/xlib/input/keyboard-encoding.html
func getKeyValue(key, state int) int {
	keyVal := keysyms[key][0]
	// keyPad
	if (keyVal >= int(sparta.KeyPadSpace)) && (keyVal <= int(sparta.KeyPad9)) {
		if (state & xgb.ModMaskShift) != 0 {
			return keyVal
		}
		return keysyms[key][1]
	}
	if (keyVal & int(sparta.KeyNoChar)) != 0 {
		if (keyVal - 1) == sparta.KeyShift {
			keyVal = int(sparta.KeyShift)
		}
		if (keyVal - 1) == sparta.KeyControl {
			keyVal = int(sparta.KeyControl)
		}
		return keyVal
	}
	if ((state & xgb.ModMaskShift) == 0) && ((state & xgb.ModMaskLock) == 0) {
		if (state & xgb.ModMask5) != 0 {
			return keysyms[key][4]
		}
		return keyVal
	}
	if ((state & xgb.ModMaskShift) == 0) && ((state & xgb.ModMaskLock) != 0) {
		if unicode.IsLetter(rune(keyVal)) {
			if unicode.IsLower(rune(keyVal)) {
				return keysyms[key][1]
			}
		}
		return keyVal
	}
	if ((state & xgb.ModMaskShift) != 0) && ((state & xgb.ModMaskLock) != 0) {
		if unicode.IsLetter(rune(keyVal)) {
			if unicode.IsLower(rune(keyVal)) {
				return keyVal
			}
		}
		return keysyms[key][1]
	}
	if (state & xgb.ModMaskShift) != 0 {
		return keysyms[key][1]
	}
	return keyVal
}
