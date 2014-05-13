// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

package sparta

// MouseButton is a mouse button.
type MouseButton int

// Mouse button values.
const (
	MouseLeft  MouseButton = 1
	MouseRight             = 2
	MouseWheel             = 3
	Mouse2                 = 4
)

// StateKey is an state key.
type StateKey int

// State key values.
const (
	StateShift   StateKey = 1
	StateLock             = 2
	StateCtrl             = 4
	StateAltGr            = 128
	StateButtonL          = 256
	StateButton2          = 512
	StateButtonR          = 1024
	StateAny              = 7 | StateAltGr | StateButtonL | StateButton2 | StateButtonR
)

// Key is a keyboard key.
type Key int

// Key values.
const (
	// flag used to recognize non-char keys
	KeyNoChar Key = 0x8000

	// General keyboard keys
	KeyBackSpace  Key = 0xff08
	KeyTab            = 0xff09
	KeyClear          = 0xff0b
	KeyReturn         = 0xff0d
	KeyPause          = 0xff13
	KeyScrollLock     = 0xff14
	KeySysReq         = 0xff15
	KeyEscape         = 0xff1b
	KeyHome           = 0xff50
	KeyLeft           = 0xff51
	KeyUp             = 0xff52
	KeyRight          = 0xff53
	KeyDown           = 0xff54
	KeyPageUp         = 0xff55
	KeyPageDown       = 0xff56
	KeyEnd            = 0xff57
	KeySelect         = 0xff60
	KeyPrint          = 0xff61
	KeyExecute        = 0xff62
	KeyInsert         = 0xff63
	KeyMenu           = 0xff67
	KeyHelp           = 0xff6a
	KeyNumLock        = 0xff7f
	KeyShift          = 0xffe1 // either one
	KeyControl        = 0xffe3 // either one
	KeyCapsLock       = 0xffe5
	KeyAlt            = 0xffe9
	KeyAltGr          = 0xffea
	KeySuperL         = 0xffeb
	KeySuperR         = 0xffec
	KeyDelete         = 0xffff

	// KeyPad values
	KeyPadSpace     Key = 0xff80
	KeyPadTab           = 0xff89
	KeyPadEnter         = 0xff8d
	KeyPadF1            = 0xff91
	KeyPadF2            = 0xff92
	KeyPadF3            = 0xff93
	KeyPadF4            = 0xff94
	KeyPadHome          = 0xff95
	KeyPadLeft          = 0xff96
	KeyPadUp            = 0xff97
	KeyPadRight         = 0xff98
	KeyPadDown          = 0xff99
	KeyPadPrior         = 0xff9a
	KeyPadPage_Up       = 0xff9a
	KeyPadNext          = 0xff9b
	KeyPadPage_Down     = 0xff9b
	KeyPadEnd           = 0xff9c
	KeyPadBegin         = 0xff9d
	KeyPadInsert        = 0xff9e
	KeyPadDelete        = 0xff9f
	KeyPadEqual         = 0xffbd
	KeyPadMultiply      = 0xffaa
	KeyPadAdd           = 0xffab
	KeyPadSeparator     = 0xffac
	KeyPadSubtract      = 0xffad
	KeyPadDecimal       = 0xffae
	KeyPadDivide        = 0xffaf
	KeyPad0             = 0xffb0
	KeyPad1             = 0xffb1
	KeyPad2             = 0xffb2
	KeyPad3             = 0xffb3
	KeyPad4             = 0xffb4
	KeyPad5             = 0xffb5
	KeyPad6             = 0xffb6
	KeyPad7             = 0xffb7
	KeyPad8             = 0xffb8
	KeyPad9             = 0xffb9

	// Function keys
	KeyF1  Key = 0xffbe
	KeyF2      = 0xffbf
	KeyF3      = 0xffc0
	KeyF4      = 0xffc1
	KeyF5      = 0xffc2
	KeyF6      = 0xffc3
	KeyF7      = 0xffc4
	KeyF8      = 0xffc5
	KeyF9      = 0xffc6
	KeyF10     = 0xffc7
	KeyF11     = 0xffc8
	KeyF12     = 0xffc9
	KeyF13     = 0xffca
	KeyF14     = 0xffcb
	KeyF15     = 0xffcc
	KeyF16     = 0xffcd
	KeyF17     = 0xffce
	KeyF18     = 0xffcf
	KeyF19     = 0xffd0
	KeyF20     = 0xffd1
	KeyF21     = 0xffd2
	KeyF22     = 0xffd3
	KeyF23     = 0xffd4
	KeyF24     = 0xffd5
)
