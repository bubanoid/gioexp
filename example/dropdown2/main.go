package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	go func() {
		//dd := property.NewDropDown([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
		w := app.NewWindow(
			app.Size(unit.Dp(900), unit.Dp(700)),
		)

		uii := New(w)

		var ops op.Ops
		for {
			switch e := w.NextEvent().(type) {
			case system.DestroyEvent:
				os.Exit(0)
				return
			case system.FrameEvent:
				gtx := l.NewContext(&ops, e)
				//th := material.NewTheme()
				//pgtx := gtx

				uii.Layout(gtx)

				//layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				//	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				//		return dd.Layout(th, pgtx, gtx)
				//	}))

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
