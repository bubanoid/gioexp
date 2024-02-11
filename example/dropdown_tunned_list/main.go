package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"

	property "github.com/arl/gioexp/component/dropdown_tunned_list"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	ui := NewUI(material.NewTheme())

	go func() {
		w := app.NewWindow(app.Title("Property List"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	th    *material.Theme
	plist *property.List
	dd    *property.DropDown

	btn widget.Clickable
}

var (
	aliceBlue = color.NRGBA{R: 240, G: 248, B: 255, A: 255}
)

func NewUI(theme *material.Theme) *UI {
	ui := &UI{
		th: theme,
		dd: property.NewDropDown([]string{"ciao", "bonjour", "hello", "hallo", "buongiorno", "buenos dias", "ola", "bom dia"}),
	}

	plist := property.NewList()

	plist.Add(ui.dd)

	ui.plist = plist
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			if e.Name == key.NameEscape {
				return nil
			}
		case system.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx C) D {
	if ui.btn.Clicked(gtx) {
		ui.dd.Selected = 2
	}

	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = 200
			return ui.plist.Layout(ui.th, gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = 200
			return material.Button(ui.th, &ui.btn, "toggle editable").Layout(gtx)
		}),
	)
}

var (
	red       = rgb(0xff0000)
	green     = rgb(0x00ff00)
	blue      = rgb(0x0000ff)
	lightGrey = rgb(0xd3d3d3)
	darkGrey  = rgb(0xa9a9a9)
)

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
