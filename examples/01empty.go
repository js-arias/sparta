// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// Sparta is a simple widget library for x11 (linux, etc.) and windows.
//
// Sparta include a sample widget set (in the directory sparta/widget),
// but it is possible to create your own custom widgets.
//
// In this simple example two empty windows will be created.

package main

// Package "init" provides the initialization of the backend. If you
// don't want to use init, or the default sparta backends, you can use
// any backend that implements the sparta.Window interface, and some
// sparta functions (such as NewWindow, Run, etc.).

import (
	"github.com/js-arias/sparta"
	_ "github.com/js-arias/sparta/init"
	"github.com/js-arias/sparta/widget"
)

func main() {
	// MainWindow is a rootless widget that can contain other widgets.
	// You can create as many as you want. You should drefine a name and
	// a title of the window (depending on the backend, this name will
	// be captured by the OS).
	//
	// Names are not necessary, but it is a good practice to have names,
	// so this will make the identification of widgets earier (using the
	// sparta.Name property).
	widget.NewMainWindow("one", "Window one")
	widget.NewMainWindow("two", "Window two")

	// sparta.Run runs the event loop of the application backend.
	sparta.Run()
}
