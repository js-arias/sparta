// Copyright (c) 2014, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in LICENSE file.

// Sparta uses properties to know information about a widget and to set
// this information. This mechanism provides a simple way to manipulate
// widgets without complex function calls.
//
// This example is two fold: first it shows the usage of properties, and
// second, it shows how to use the canvas widget.
//
// The canvas widget is a widget in which objects and text can be draw at
// will. In most other objects, the client code is not expected to do
// any drawing.
//
// In this example, the main window will be split in four parts (childs),
// each one, a canvas, will draw different things: text, a sine function
// (a set of lines), two polygons, and a set of geometrical objects. Most
// of the example is based on C. Petzold (1998) "Programming Windows"
// Microsoft Press. (figs. 5.6, 5.14, and 5.21).

package main

import (
	"image"
	"image/color"
	"math"

	"github.com/js-arias/sparta"
	_ "github.com/js-arias/sparta/init"
	"github.com/js-arias/sparta/widget"
)

// polygon reference points
var polyPts = []image.Point{
	image.Pt(10, 70),
	image.Pt(50, 70),
	image.Pt(50, 10),
	image.Pt(90, 10),
	image.Pt(90, 50),
	image.Pt(30, 50),
	image.Pt(30, 90),
	image.Pt(70, 90),
	image.Pt(70, 30),
	image.Pt(10, 30),
}

// objects data
type objData struct {
	l1, l2 []image.Point
	arc    widget.Arc
}

// page data (used for the text display)
type pageData struct {
	pos  int
	page int
}

func main() {
	m := widget.NewMainWindow("main", "Canvas")

	// The method Property returns an interface of the asked property. If
	// the widget does not process that property, it will return nil.
	// sparta.Geometry is a property that indicate the dimensions (and
	// position) of the widget relative to its parent.
	geo := m.Property(sparta.Geometry).(image.Rectangle)

	// NewCanvas creates a new drawable canvas.
	sn := widget.NewCanvas(m, "sine", image.Rect(0, 0, geo.Dx()/2, geo.Dy()/2))
	// The method SetProperty sets a property in the widget.
	// We want that the widget have a border, so we set sparta.Border as true.
	sn.SetProperty(sparta.Border, true)
	// create the points for the sine function.
	pts := make([]image.Point, 1000)
	// set the initial values of the sine function.
	sine(pts, geo.Dx()/2, geo.Dy()/2)
	// sparta.Data is a property that store client defined data of a widget.
	sn.SetProperty(sparta.Data, pts)
	// We want to process the configure event, to update the position of
	// the sine points.
	sn.Capture(sparta.Configure, snConf)
	// We wnat to process the expose event, to draw the content of the widget.
	sn.Capture(sparta.Expose, snExpose)
	// We send and update request to guarantee that the content of the
	// widget will be drawn.
	sn.Update()

	pg := widget.NewCanvas(m, "polygon", image.Rect(geo.Dx()/2, 0, geo.Dx(), geo.Dy()/2))
	pg.SetProperty(sparta.Border, true)

	// We set the background of the widget to be light grey, using the
	// property sparta.Background.
	pg.SetProperty(sparta.Background, color.RGBA{190, 190, 190, 0})
	// We set the foreground of the widget (i.e. the color used to draw
	// objects to be drak grey, using the property sparta.Foreground.
	pg.SetProperty(sparta.Foreground, color.RGBA{90, 90, 90, 0})

	poly := []widget.Polygon{
		widget.Polygon{Pt: make([]image.Point, len(polyPts))},
		widget.Polygon{Pt: make([]image.Point, len(polyPts)), Fill: true},
	}
	setPoly(poly, geo.Dx()-(geo.Dx()/2), geo.Dy()/2)
	pg.SetProperty(sparta.Data, poly)
	pg.Capture(sparta.Configure, pgConf)
	pg.Capture(sparta.Expose, pgExpose)
	pg.Update()

	ob := widget.NewCanvas(m, "objects", image.Rect(0, geo.Dy()/2, geo.Dx()/2, geo.Dy()))
	ob.SetProperty(sparta.Border, true)
	objs := &objData{
		l1: make([]image.Point, 2),
		l2: make([]image.Point, 2),
		arc: widget.Arc{
			Angle2: math.Pi * 2,
			Fill:   true,
		},
	}
	setObj(objs, geo.Dx()/2, geo.Dy()-(geo.Dy()/2))
	ob.SetProperty(sparta.Data, objs)
	ob.Capture(sparta.Configure, obConf)
	ob.Capture(sparta.Expose, obExpose)
	ob.Update()

	tx := widget.NewCanvas(m, "poem", image.Rect(geo.Dx()/2, geo.Dy()/2, geo.Dx(), geo.Dy()))
	tx.SetProperty(sparta.Border, true)
	// In sparta the size of a text glyph is given by sparta.HeightUnit
	// and sparta.WidthUnit.
	txtData := &pageData{
		pos:  0,
		page: (geo.Dy() - (geo.Dy() / 2)) / sparta.HeightUnit,
	}
	tx.SetProperty(sparta.Data, txtData)
	tx.Capture(sparta.Configure, txConf)
	tx.Capture(sparta.Expose, txExpose)

	// To navigate the text we also capture the mouse and the keyboard.
	tx.Capture(sparta.Mouse, txMouse)
	tx.Capture(sparta.KeyEv, txKey)
	tx.Update()

	// We want to capture the mainwindow configure event so if the
	// mainWindow changes its size, we can update the size and position
	// of the all other widgets.
	m.Capture(sparta.Configure, mConf)

	sparta.Run()
}

func mConf(m sparta.Widget, e interface{}) bool {
	// the sparta.Childs property return the children of a widget.
	ch := m.Property(sparta.Childs).([]sparta.Widget)
	ev := e.(sparta.ConfigureEvent)
	for _, c := range ch {
		// We check the name propery of each widget and use it
		// to set the new geometry of each widget.
		switch nm := c.Property(sparta.Name).(string); nm {
		case "sine":
			c.SetProperty(sparta.Geometry, image.Rect(0, 0, ev.Rect.Dx()/2, ev.Rect.Dy()/2))
		case "polygon":
			c.SetProperty(sparta.Geometry, image.Rect(ev.Rect.Dx()/2, 0, ev.Rect.Dx(), ev.Rect.Dy()/2))
		case "objects":
			c.SetProperty(sparta.Geometry, image.Rect(0, ev.Rect.Dy()/2, ev.Rect.Dx()/2, ev.Rect.Dy()))
		case "poem":
			c.SetProperty(sparta.Geometry, image.Rect(ev.Rect.Dx()/2, ev.Rect.Dy()/2, ev.Rect.Dx(), ev.Rect.Dy()))
		}
	}
	return false
}

// SnConf sets the new values of the sine function points
func snConf(sn sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ConfigureEvent)

	// Get data from the widget.
	data := sn.Property(sparta.Data).([]image.Point)
	sine(data, ev.Rect.Dx(), ev.Rect.Dy())
	return false
}

// Sine set values for the sine function.
func sine(pts []image.Point, width, height int) {
	for i := 0; i < len(pts); i++ {
		pts[i].X = i * width / len(pts)
		pts[i].Y = int(float64(height) * (1 - math.Sin(2*math.Pi*float64(i)/float64(len(pts)))) / 2)
	}
}

// SnExpose draw the sine function.
func snExpose(sn sparta.Widget, e interface{}) bool {
	// get point data.
	data := sn.Property(sparta.Data).([]image.Point)
	// get widget geometry
	geo := sn.Property(sparta.Geometry).(image.Rectangle)

	c := sn.(*widget.Canvas)
	// The widget canvas ruses the function Draw to draw particular
	// objects, it depends on the data type to decide what to draw.
	// Here a line (the "x" axis), is draw.
	c.Draw([]image.Point{image.Pt(0, geo.Dy()/2), image.Pt(geo.Dx(), geo.Dy()/2)})
	// Then the sine function is draw.
	c.Draw(data)

	return false
}

// PgConf sets the new values of the polygon points.
func pgConf(pg sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ConfigureEvent)
	data := pg.Property(sparta.Data).([]widget.Polygon)
	setPoly(data, ev.Rect.Dx(), ev.Rect.Dy())
	return false
}

// SetPoly sets the polygon values.
func setPoly(poly []widget.Polygon, width, height int) {
	for i, p := range polyPts {
		poly[0].Pt[i].X = width * p.X / 200
		poly[1].Pt[i].X = poly[0].Pt[i].X + (width / 2)
		poly[0].Pt[i].Y = height * p.Y / 100
		poly[1].Pt[i].Y = poly[0].Pt[i].Y
	}
}

// pgExpose draws the polygons.
func pgExpose(pg sparta.Widget, e interface{}) bool {
	data := pg.Property(sparta.Data).([]widget.Polygon)
	c := pg.(*widget.Canvas)
	c.Draw(data[0])
	c.Draw(data[1])
	return false
}

// obConf sets the new values of the objects points.
func obConf(ob sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ConfigureEvent)
	data := ob.Property(sparta.Data).(*objData)
	setObj(data, ev.Rect.Dx(), ev.Rect.Dy())
	return false
}

// SetObj sets the object values.
func setObj(objs *objData, width, height int) {
	objs.arc.Rect = image.Rect(width/8, height/8, 7*width/8, 7*height/8)
	objs.l1[1].X, objs.l1[1].Y = width, height
	objs.l2[0].X, objs.l2[1].Y = width, height
}

// ObExpose draws the objects.
func obExpose(ob sparta.Widget, e interface{}) bool {
	data := ob.Property(sparta.Data).(*objData)
	c := ob.(*widget.Canvas)
	// Set color sets a color temporalely for the following dawing
	// operations. Contrast this with the property sparta.Background
	// and sparta.Foreground.
	c.SetColor(sparta.Foreground, color.RGBA{255, 0, 0, 0})
	c.Draw(widget.Rectangle{Rect: data.arc.Rect})
	c.SetColor(sparta.Foreground, color.RGBA{0, 255, 0, 0})
	c.Draw(data.l1)
	c.Draw(data.l2)
	c.SetColor(sparta.Foreground, color.RGBA{0, 0, 255, 0})
	c.Draw(data.arc)
	return false
}

// TxConf sets the txt data values.
func txConf(tx sparta.Widget, e interface{}) bool {
	ev := e.(sparta.ConfigureEvent)
	data := tx.Property(sparta.Data).(*pageData)
	data.page = ev.Rect.Dy() / sparta.HeightUnit
	return false
}

// TxExpose draws the poem.
func txExpose(tx sparta.Widget, e interface{}) bool {
	data := tx.Property(sparta.Data).(*pageData)
	rect := tx.Property(sparta.Geometry).(image.Rectangle)
	c := tx.(*widget.Canvas)

	// Text store the text to be drawing
	txt := widget.Text{}
	txt.Pos.X = 2
	for i, ln := range poem[data.pos:] {
		// The position of the text is the top-right corner of
		// the rectange that enclose the text.
		txt.Pos.Y = (i * sparta.HeightUnit) + 2
		if txt.Pos.Y > rect.Dy() {
			break
		}
		txt.Text = ln
		c.Draw(txt)
	}
	return false
}

// TxKey gets keyboard events.
func txKey(tx sparta.Widget, e interface{}) bool {
	data := tx.Property(sparta.Data).(*pageData)
	ev := e.(sparta.KeyEvent)
	switch ev.Key {
	case sparta.KeyDown:
		if (data.pos + 1) < (len(poem) - data.page + 1) {
			data.pos++
		}
		tx.Update()
	case sparta.KeyUp:
		if (data.pos - 1) >= 0 {
			data.pos--
		}
		tx.Update()
	case sparta.KeyPageUp:
		if data.pos == 0 {
			break
		}
		data.pos -= data.page
		if data.pos < 0 {
			data.pos = 0
		}
		tx.Update()
	case sparta.KeyPageDown:
		if data.pos == (len(poem) - data.page) {
			break
		}
		data.pos += data.page
		if data.pos > (len(poem) - data.page + 1) {
			data.pos = len(poem) - data.page
		}
		tx.Update()
	}
	return true
}

// TxMouse gets mouse events.
func txMouse(tx sparta.Widget, e interface{}) bool {
	data := tx.Property(sparta.Data).(*pageData)
	ev := e.(sparta.MouseEvent)
	switch ev.Button {
	case sparta.MouseLeft, -sparta.MouseWheel:
		if (data.pos + 1) < (len(poem) - data.page + 1) {
			data.pos++
		}
		tx.Update()
	case sparta.MouseRight, sparta.MouseWheel:
		if (data.pos - 1) >= 0 {
			data.pos--
		}
		tx.Update()
	}
	return true
}

var poem = []string{
	"MY good blade carves the casques of men,",
	"My tough lance thrusteth sure,",
	"My strength is as the strength of ten,",
	"Because my heart is pure.",
	"The shattering trumpet shrilleth high,",
	"The hard brands shiver on the steel,",
	"The splinter'd spear-shafts crack and fly,",
	"The horse and rider reel:",
	"They reel, they roll in clanging lists,",
	"And when the tide of combat stands,",
	"Perfume and flowers fall in showers,",
	"That lightly rain from ladies' hands.",
	"How sweet are looks that ladies bend",
	"On whom their favours fall !",
	"For them I battle till the end,",
	"To save from shame and thrall:",
	"But all my heart is drawn above,",
	"My knees are bow'd in crypt and shrine:",
	"I never felt the kiss of love,",
	"Nor maiden's hand in mine.",
	"More bounteous aspects on me beam,",
	"Me mightier transports move and thrill;",
	"So keep I fair thro' faith and prayer",
	"A virgin heart in work and will.",
	"When down the stormy crescent goes,",
	"A light before me swims,",
	"Between dark stems the forest glows,",
	"I hear a noise of hymns:",
	"Then by some secret shrine I ride;",
	"I hear a voice but none are there;",
	"The stalls are void, the doors are wide,",
	"The tapers burning fair.",
	"Fair gleams the snowy altar-cloth,",
	"The silver vessels sparkle clean,",
	"The shrill bell rings, the censer swings,",
	"And solemn chaunts resound between.",
	"Sometimes on lonely mountain-meres",
	"I find a magic bark;",
	"I leap on board: no helmsman steers:",
	"I float till all is dark.",
	"A gentle sound, an awful light !",
	"Three arngels bear the holy Grail:",
	"With folded feet, in stoles of white,",
	"On sleeping wings they sail.",
	"Ah, blessed vision! blood of God!",
	"AsMy spirit beats her mortal bars,",
	"down dark tides the glory slides,",
	"And star-like mingles with the stars.",
	"When on my goodly charger borne",
	"Thro' dreaming towns I go,",
	"The cock crows ere the Christmas morn,",
	"The streets are dumb with snow.",
	"The tempest crackles on the leads,",
	"And, ringing, springs from brand and mail;",
	"But o'er the dark a glory spreads,",
	"And gilds the driving hail.",
	"I leave the plain, I climb the height;",
	"No branchy thicket shelter yields;",
	"But blessed forms in whistling storms",
	"Fly o'er waste fens and windy fields.",
	"A maiden knight--to me is given",
	"Such hope, I know not fear;",
	"I yearn to breathe the airs of heaven",
	"That often meet me here.",
	"I muse on joy that will not cease,",
	"Pure spaces clothed in living beams,",
	"Pure lilies of eternal peace,",
	"Whose odours haunt my dreams;",
	"And, stricken by an angel's hand,",
	"This mortal armour that I wear,",
	"This weight and size, this heart and eyes,",
	"Are touch'd, are turn'd to finest air.",
	"The clouds are broken in the sky,",
	"And thro' the mountain-walls",
	"A rolling organ-harmony",
	"Swells up, and shakes and falls.",
	"Then move the trees, the copses nod,",
	"Wings flutter, voices hover clear:",
	"'O just and faithful knight of God!",
	"Ride on ! the prize is near.'",
	"So pass I hostel, hall, and grange;",
	"By bridge and ford, by park and pale,",
	"All-arm'd I ride, whate'er betide,",
	"Until I find the holy Grail.",
}
