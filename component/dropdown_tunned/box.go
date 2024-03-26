package dropdown_tunned

import (
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

// todo (AA): usage
//   dropdown_tunned.StrokeRect(gtx, gtx.Ops, res)

func StrokeRect(gtx C, ops *op.Ops, min D, delta D) {
	rect := clip.Rect{
		Min: image.Point{X: min.Size.X, Y: min.Size.Y},
		Max: image.Point{X: min.Size.X + delta.Size.X, Y: min.Size.Y + delta.Size.Y},
	}
	paint.FillShape(ops, color.NRGBA{R: 0x80, A: 0xFF},
		clip.Stroke{
			Path:  rect.Path(),
			Width: 4,
		}.Op(),
	)
}
