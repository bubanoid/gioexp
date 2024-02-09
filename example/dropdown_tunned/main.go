package main

import (
	"github.com/arl/gioexp/component/dropdown_tunned"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	ui := NewUI(material.NewTheme())

	go func() {
		w := app.NewWindow(app.Title("Property DropDown"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	th  *material.Theme
	dd  *dropdown_tunned.DropDown
	btn widget.Clickable
}

func NewUI(theme *material.Theme) *UI {
	ui := &UI{
		th: theme,
	}

	dropdown := dropdown_tunned.NewDropdown(
		[]string{"ciao", "bonjour", "hello", "hallo", "buongiorno", "buenos dias", "ola", "bom dia"},
	)

	ui.dd = dropdown
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
		ui.dd.DdWidget.Selected = 2
	}

	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return ui.dd.Layout(ui.th, gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return material.Button(ui.th, &ui.btn, "toggle editable").Layout(gtx)
				}),
			)
		}),
	)
}
