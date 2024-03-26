package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"os"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	th := material.NewTheme()

	go func() {
		w := app.NewWindow(app.Title("DropDown Menu"))
		for {
			switch e := w.NextEvent().(type) {
			case system.FrameEvent:
				//gtx := layout.NewContext(new(op.Ops), e)
				gtx := layout.Context{
					Ops: new(op.Ops),
					Constraints: layout.Constraints{
						Max: image.Point{X: 100, Y: 100},
					},
				}

				res := layout.Center.Layout(gtx, func(gtx C) D {
					stack := layout.Stack{}
					off1 := op.Offset(image.Point{X: 0, Y: 10}).Push(gtx.Ops)
					res := stack.Layout(gtx,
						// Widget 1, placed at the bottom
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							off2 := op.Offset(image.Point{X: 0, Y: 0}).Push(gtx.Ops)
							res := material.Button(th, new(widget.Clickable), "Button 1").Layout(gtx)
							off2.Pop()
							return res
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							off2 := op.Offset(image.Point{X: 10, Y: -10}).Push(gtx.Ops)
							res := material.Label(th, unit.Sp(12), "Button 2").Layout(gtx)
							off2.Pop()
							//fmt.Println(res)
							return res
						}),

						// Widget 2, placed on top of Widget 1
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							off2 := op.Offset(image.Point{X: 10, Y: -10}).Push(gtx.Ops)
							res := material.Label(th, unit.Sp(12), "Button 2").Layout(gtx)
							off2.Pop()
							//fmt.Println(res)
							return res
						}),
					)
					off1.Pop()
					fmt.Println(res)
					return res
				})
				fmt.Println(res)
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()

	//gtx := layout.Context{
	//	Ops: new(op.Ops),
	//	Constraints: layout.Constraints{
	//		Max: image.Point{X: 100, Y: 100},
	//	},
	//}
	//
	//stack := layout.Stack{}
	//stack.Layout(gtx,
	//	// Widget 1, placed at the bottom
	//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//		return material.Button(th, new(widget.Clickable), "Button 1").Layout(gtx)
	//	}),
	//	// Widget 2, placed on top of Widget 1
	//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//		return material.Button(th, new(widget.Clickable), "Button 2").Layout(gtx)
	//	}),
	//)

	//layout.Stack{}.Layout(gtx,
	//	// Force widget to the same size as the second.
	//	layout.Expanded(func(gtx layout.Context) layout.Dimensions {
	//		fmt.Printf("Expand: %v\n", gtx.Constraints)
	//		return layoutWidget(gtx, 10, 10)
	//	}),
	//	// Rigid 50x50 widget.
	//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//		return layoutWidget(gtx, 50, 50)
	//	}),
	//)

}

func layoutWidget(ctx layout.Context, width, height int) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: width,
			Y: height,
		},
	}
}
