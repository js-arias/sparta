// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// +build windows

package win32

import (
	"image"
	"log"
	"os"
	"unicode/utf16"

	"github.com/AllenDang/w32"
	"github.com/js-arias/sparta"
)

func init() {
	sparta.Run = run
	sparta.Close = closeApp
	sparta.SendEvent = sendEvent
}

// SendEvent sends an event to a widget.
func sendEvent(dest sparta.Widget, comm sparta.CommandEvent) {
	dwin := dest.Window().(*window)
	var id w32.HWND
	if comm.Source != nil {
		id = comm.Source.Window().(*window).id
	}
	w32.PostMessage(dwin.id, w32.WM_USER, uintptr(id), uintptr(comm.Value))
}

func run() {
	msg := &w32.MSG{}
	for {
		switch val := w32.GetMessage(msg, 0, 0, 0); val {
		case -1:
			log.Printf("w32: error: %v\n", getLastError())
			os.Exit(1)
		case 0:
			return
		default:
			w32.TranslateMessage(msg)
			w32.DispatchMessage(msg)
		}
	}
}

// CloseApp closes the application.
func closeApp() {
	for _, b := range mapBrush {
		w32.DeleteObject(w32.HGDIOBJ(b.pen))
		w32.DeleteObject(w32.HGDIOBJ(b.brush))
	}
	w32.PostQuitMessage(0)
}

// WinEvent proccess a win32 event.
func winEvent(id w32.HWND, event uint32, wParam, lParam uintptr) uintptr {
	w, ok := widgetTable[id]
	if !ok {
		return w32.DefWindowProc(id, event, wParam, lParam)
	}
	switch event {
	case w32.WM_CHAR:
		r := utf16.Decode([]uint16{uint16(loWord(uint32(wParam)))})
		if len(r) == 0 {
			break
		}
		key := r[0]
		if key == 0 {
			break
		}
		if (key == '\b') || (key == '\t') || (key == '\n') || (key == '\r') {
			break
		}
		ev := sparta.KeyEvent{
			Key:   sparta.Key(key),
			State: getState(),
		}
		x, y, _ := w32.GetCursorPos()
		ev.Loc.X, ev.Loc.Y, _ = w32.ScreenToClient(id, x, y)
		w.OnEvent(ev)
	case w32.WM_CLOSE:
		if w.Property(sparta.Parent) != nil {
			break
		}
		w.OnEvent(sparta.CloseEvent{})
	case w32.WM_KEYDOWN:
		key := getKeyValue(wParam)
		if key == 0 {
			break
		}
		if (key & sparta.KeyNoChar) == 0 {
			break
		}
		ev := sparta.KeyEvent{
			Key:   key,
			State: getState(),
		}
		x, y, _ := w32.GetCursorPos()
		ev.Loc.X, ev.Loc.Y, _ = w32.ScreenToClient(id, x, y)
		w.OnEvent(ev)
	case w32.WM_KEYUP:
		key := getKeyValue(wParam)
		if key == 0 {
			break
		}
		ev := sparta.KeyEvent{
			Key:   -key,
			State: getState(),
		}
		x, y, _ := w32.GetCursorPos()
		ev.Loc.X, ev.Loc.Y, _ = w32.ScreenToClient(id, x, y)
		w.OnEvent(ev)
	case w32.WM_LBUTTONDOWN, w32.WM_RBUTTONDOWN, w32.WM_MBUTTONDOWN:
		ev := sparta.MouseEvent{
			Button: getButton(event),
			State:  getState(),
			Loc:    image.Pt(getXLParam(lParam), getYLParam(lParam)),
		}
		w.OnEvent(ev)
		w.Focus()
	case w32.WM_LBUTTONUP, w32.WM_RBUTTONUP, w32.WM_MBUTTONUP:
		ev := sparta.MouseEvent{
			Button: -getButton(event),
			State:  getState(),
			Loc:    image.Pt(getXLParam(lParam), getYLParam(lParam)),
		}
		w.OnEvent(ev)
	case w32.WM_MOUSEMOVE:
		ev := sparta.MouseEvent{
			Loc: image.Pt(getXLParam(lParam), getYLParam(lParam)),
		}
		w.OnEvent(ev)
	case w32.WM_MOUSEWHEEL:
		ev := sparta.MouseEvent{
			Button: sparta.MouseWheel,
		}
		if getWheelDeltaWParam(wParam) < 0 {
			ev.Button = -sparta.MouseWheel
		}
		ev.Loc.X, ev.Loc.Y, _ = w32.ScreenToClient(id, getXLParam(lParam), getYLParam(lParam))
		w = propagateWheel(w, ev.Loc)
		w.OnEvent(ev)
	case w32.WM_MOVE:
		win := w.Window().(*window)
		win.pos.X, win.pos.Y = int(loWord(uint32(lParam))), int(hiWord(uint32(lParam)))
	case w32.WM_PAINT:
		win := w.Window().(*window)
		ps := &w32.PAINTSTRUCT{}
		win.dc = w32.BeginPaint(id, ps)
		win.isPaint = true

		w32.SetBkMode(win.dc, w32.TRANSPARENT)
		w32.SetBkColor(win.dc, win.back.color)

		// "clear" the area
		w32.SelectObject(win.dc, w32.HGDIOBJ(win.back.brush))
		w32.SelectObject(win.dc, w32.HGDIOBJ(win.back.pen))
		w32.Rectangle(win.dc, int(ps.RcPaint.Left), int(ps.RcPaint.Top), int(ps.RcPaint.Right), int(ps.RcPaint.Bottom))

		w32.SelectObject(win.dc, w32.HGDIOBJ(win.fore.brush))
		w32.SelectObject(win.dc, w32.HGDIOBJ(win.fore.pen))
		w32.SelectObject(win.dc, w32.HGDIOBJ(winFont))
		w32.SetTextColor(win.dc, win.fore.color)
		win.curr = win.fore

		ev := sparta.ExposeEvent{image.Rect(int(ps.RcPaint.Left), int(ps.RcPaint.Top), int(ps.RcPaint.Right), int(ps.RcPaint.Bottom))}
		w.OnEvent(ev)
		w32.EndPaint(id, ps)
		win.isPaint = false
		win.dc = 0
	case w32.WM_SIZE:
		win := w.Window().(*window)
		ev := sparta.ConfigureEvent{image.Rect(win.pos.X, win.pos.Y, win.pos.X+int(loWord(uint32(lParam))), win.pos.Y+int(hiWord(uint32(lParam))))}
		w.OnEvent(ev)
	case w32.WM_USER:
		src, ok := widgetTable[w32.HWND(wParam)]
		if !ok {
			src = nil
		}
		ev := sparta.CommandEvent{
			Source: src,
			Value:  int(int32(lParam)),
		}
		w.OnEvent(ev)
	default:
		return w32.DefWindowProc(id, event, wParam, lParam)
	}
	return 0
}

func getButton(event uint32) sparta.MouseButton {
	switch event {
	case w32.WM_LBUTTONDOWN, w32.WM_LBUTTONUP:
		return sparta.MouseLeft
	case w32.WM_RBUTTONDOWN, w32.WM_RBUTTONUP:
		return sparta.MouseRight
	case w32.WM_MBUTTONDOWN, w32.WM_MBUTTONUP:
		return sparta.Mouse2
	}
	return 0
}

func propagateWheel(w sparta.Widget, pt image.Point) sparta.Widget {
	rect := w.Property(sparta.Geometry).(image.Rectangle)
	childs := w.Property(sparta.Childs)
	if childs == nil {
		return w
	}
	for _, ch := range childs.([]sparta.Widget) {
		rect = ch.Property(sparta.Geometry).(image.Rectangle)
		if pt.In(rect) {
			return propagateWheel(ch, pt.Add(rect.Min))
		}
	}
	return w
}

func getState() sparta.StateKey {
	var state sparta.StateKey
	if GetKeyState(w32.VK_RMENU) < 0 {
		state = sparta.StateAltGr
	}
	if GetKeyState(w32.VK_SHIFT) < 0 {
		state |= sparta.StateShift
	}
	if GetKeyState(w32.VK_CONTROL) < 0 {
		state |= sparta.StateCtrl
	}
	if GetKeyState(w32.VK_LBUTTON) < 0 {
		state |= sparta.StateButtonL
	}
	if GetKeyState(w32.VK_MBUTTON) < 0 {
		state |= sparta.StateButton2
	}
	if GetKeyState(w32.VK_RBUTTON) < 0 {
		state |= sparta.StateButtonR
	}
	if (GetKeyState(w32.VK_CAPITAL) & 1) != 0 {
		state |= sparta.StateLock
	}
	return state
}

func getKeyValue(wParam uintptr) sparta.Key {
	switch wParam {
	case w32.VK_BACK:
		return sparta.KeyBackSpace
	case w32.VK_TAB:
		return sparta.KeyTab
	case w32.VK_CLEAR:
		return sparta.KeyClear
	case w32.VK_RETURN:
		return sparta.KeyReturn
	case w32.VK_SHIFT:
		return sparta.KeyShift
	case w32.VK_CONTROL:
		return sparta.KeyControl
	case w32.VK_MENU:
		return sparta.KeyMenu
	case w32.VK_PAUSE:
		return sparta.KeyPause
	case w32.VK_CAPITAL:
		return sparta.KeyCapsLock
	case w32.VK_ESCAPE:
		return sparta.KeyEscape
	case w32.VK_SPACE:
		return ' '
	case w32.VK_PRIOR:
		return sparta.KeyPageUp
	case w32.VK_NEXT:
		return sparta.KeyPageDown
	case w32.VK_END:
		return sparta.KeyEnd
	case w32.VK_HOME:
		return sparta.KeyHome
	case w32.VK_LEFT:
		return sparta.KeyLeft
	case w32.VK_UP:
		return sparta.KeyUp
	case w32.VK_RIGHT:
		return sparta.KeyRight
	case w32.VK_DOWN:
		return sparta.KeyDown
	case w32.VK_SELECT:
		return sparta.KeySelect
	case w32.VK_PRINT:
		return sparta.KeyPrint
	case w32.VK_EXECUTE:
		return sparta.KeyExecute
	case w32.VK_SNAPSHOT:
		return sparta.KeySysReq
	case w32.VK_INSERT:
		return sparta.KeyInsert
	case w32.VK_DELETE:
		return sparta.KeyDelete
	case w32.VK_HELP:
		return sparta.KeyHelp
	case w32.VK_LWIN:
		return sparta.KeySuperL
	case w32.VK_RWIN:
		return sparta.KeySuperR
	case w32.VK_NUMPAD0:
		return sparta.KeyPad0
	case w32.VK_NUMPAD1:
		return sparta.KeyPad1
	case w32.VK_NUMPAD2:
		return sparta.KeyPad2
	case w32.VK_NUMPAD3:
		return sparta.KeyPad3
	case w32.VK_NUMPAD4:
		return sparta.KeyPad4
	case w32.VK_NUMPAD5:
		return sparta.KeyPad5
	case w32.VK_NUMPAD6:
		return sparta.KeyPad6
	case w32.VK_NUMPAD7:
		return sparta.KeyPad7
	case w32.VK_NUMPAD8:
		return sparta.KeyPad8
	case w32.VK_NUMPAD9:
		return sparta.KeyPad9
	case w32.VK_MULTIPLY:
		return sparta.KeyPadMultiply
	case w32.VK_ADD:
		return sparta.KeyPadAdd
	case w32.VK_SEPARATOR:
		return sparta.KeyPadSeparator
	case w32.VK_SUBTRACT:
		return sparta.KeyPadSubtract
	case w32.VK_DECIMAL:
		return sparta.KeyPadDecimal
	case w32.VK_DIVIDE:
		return sparta.KeyPadDivide
	case w32.VK_F1:
		return sparta.KeyF1
	case w32.VK_F2:
		return sparta.KeyF2
	case w32.VK_F3:
		return sparta.KeyF3
	case w32.VK_F4:
		return sparta.KeyF4
	case w32.VK_F5:
		return sparta.KeyF5
	case w32.VK_F6:
		return sparta.KeyF6
	case w32.VK_F7:
		return sparta.KeyF7
	case w32.VK_F8:
		return sparta.KeyF8
	case w32.VK_F9:
		return sparta.KeyF9
	case w32.VK_F10:
		return sparta.KeyF10
	case w32.VK_F11:
		return sparta.KeyF11
	case w32.VK_F12:
		return sparta.KeyF12
	case w32.VK_F13:
		return sparta.KeyF13
	case w32.VK_F14:
		return sparta.KeyF14
	case w32.VK_F15:
		return sparta.KeyF15
	case w32.VK_F16:
		return sparta.KeyF16
	case w32.VK_F17:
		return sparta.KeyF17
	case w32.VK_F18:
		return sparta.KeyF18
	case w32.VK_F19:
		return sparta.KeyF19
	case w32.VK_F20:
		return sparta.KeyF20
	case w32.VK_F21:
		return sparta.KeyF21
	case w32.VK_F22:
		return sparta.KeyF22
	case w32.VK_F23:
		return sparta.KeyF23
	case w32.VK_F24:
		return sparta.KeyF24
	case w32.VK_NUMLOCK:
		return sparta.KeyNumLock
	case w32.VK_SCROLL:
		return sparta.KeyScrollLock
	default:
		if (wParam >= 0x30) && (wParam <= 0x39) {
			return sparta.Key(wParam)
		}
		if (wParam >= 0x41) && (wParam <= 0x5A) {
			return sparta.Key(wParam)
		}
	}
	return 0
}
