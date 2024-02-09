package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/arl/gioexp/component/property"
)

func main() {
	go func() {
		//events := make(chan event.Event)
		dd := property.NewDropDown([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
		w := app.NewWindow()
		var ops op.Ops
		for {
			//events <- w.NextEvent()
			//ee := <-events
			switch e := w.NextEvent().(type) {
			case system.DestroyEvent:
				os.Exit(0)
				return
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				th := material.NewTheme()
				pgtx := gtx
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return dd.Layout(th, pgtx, gtx)
					}))
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
