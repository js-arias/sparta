Sparta
======

Sparta is a generic (and spartan) widget toolkit for window based guis. Current 
implementation has been tested in linux (32 and 64 bits) as well as in windows 
(32 bits).

It is a "classic" toolkit in the sense that uses the traditional system of 
creating widget and answering events through an event function.

Setup
-----

    go get github.com/js-arias/sparta

In the main package initialize the init package that will authomatically 
setup the corresponding backend.

After all widgets are defined, the main loop of the program is executed using 
the Run() function.

The package defines a basic widget interface, some simple widgets are included 
in the widget package that can be used for applications, or can be used as an 
example of how the widgets can be implemented.

Properties
----------

To kept the API small, instead of a lot of function calls, each widget has a 
set of properties, some are common to all or most all widgets, and others just 
limited to particular widgets. The properties can be set using SetProperty and 
Property functions.

Events
------

It is possible to define a function based on an event. The basic events are 
produced from the Mouse, the Keyboard, Exposition and Configuration of the 
widget, as well as from other widgets or the main process.

Examples
--------

The directory example provide a collection of diferent examples that shows the 
usage of the library.

Authorship and license
----------------------

Copyright (c) 2013, J. Salvador Arias <jsalarias@gmail.com>
All rights reserved.
Distributed under BSD2 license that can be found in the LICENSE file.

