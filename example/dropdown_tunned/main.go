package main

import (
	"fmt"
	"github.com/arl/gioexp/component/dropdown_tunned"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	ui := NewUI(material.NewTheme())

	go func() {
		w := app.NewWindow(app.Title("DropDown Menu"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	th *material.Theme
	dd *dropdown_tunned.DropDown
}

func doSmth() {
	fmt.Print("dropdown item selected\n")
}

func NewUI(theme *material.Theme) *UI {
	ui := &UI{
		th: theme,
	}

	dropdown := dropdown_tunned.NewDropdown(
		[]string{"ciao", "bonjour", "hello", "hallo", "buongiorno", "buenos dias", "ola", "bom dia"},
		doSmth,
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
			layout.Center.Layout(gtx, func(gtx C) D {
				return ui.Layout(gtx)
			})
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	return ui.dd.Layout(ui.th, gtx)
}
