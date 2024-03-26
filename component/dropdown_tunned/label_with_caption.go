package dropdown_tunned

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type LabelBordered struct {
	hint   Hint
	border border
}

type Hint struct {
	TextSize   unit.Sp
	Inset      layout.Inset
	Dimensions layout.Dimensions
}

type border struct {
	Thickness unit.Dp
	Color     color.NRGBA
}

func (label *LabelBordered) Update(gtx C, th *material.Theme, hint string) {
	var (
		// Text size transitions.
		textSmall = th.TextSize * 0.8
		// Border color transitions.
		//borderColorHovered = component.WithAlpha(th.Palette.Fg, 221)
		borderColor     = th.Palette.ContrastBg // todo (AA): set border color
		borderThickness = unit.Dp(2.0)          // todo (AA): set border thickness
	)
	label.hint.TextSize = textSmall
	label.border.Thickness = borderThickness
	label.border.Color = borderColor

	// Calculate the dimensions of the hint caption size and store the
	// result for use label clipping.
	// Hack: Reset min constraint to 0 to avoid min == max.
	gtx.Constraints.Min.X = 0
	gtx.Constraints.Min.Y = 0

	macro := op.Record(gtx.Ops)
	var spacing unit.Dp
	if len(hint) > 0 {
		spacing = 4 // todo (AA): make configurable
	}
	label.hint.Dimensions = layout.Inset{
		Left:  spacing,
		Right: spacing,
	}.Layout(gtx, func(gtx C) D {
		return material.Label(th, textSmall, hint).Layout(gtx)
	})
	macro.Stop()

	labelTopInsetNormal := float32(label.hint.Dimensions.Size.Y) - float32(label.hint.Dimensions.Size.Y/4)
	topInsetDP := unit.Dp(labelTopInsetNormal / gtx.Metric.PxPerDp)
	topInsetActiveDP := (topInsetDP / 2 * -1) - label.border.Thickness
	label.hint.Inset = layout.Inset{
		Top:  topInsetActiveDP,
		Left: unit.Dp(10), // todo (AA): make configurable
	}
}

func (label *LabelBordered) Layout(gtx C, th *material.Theme, labelText string, hint string) D {
	label.Update(gtx, th, hint)
	// Offset accounts for hint height, which sticks above the border dimensions.
	defer op.Offset(image.Pt(0, label.hint.Dimensions.Size.Y/2)).Push(gtx.Ops).Pop()

	label.hint.Inset.Layout(
		gtx,
		func(gtx C) D {
			return layout.Inset{
				Left:  unit.Dp(4),
				Right: unit.Dp(4),
			}.Layout(gtx, func(gtx C) D {
				lb := material.Label(th, label.hint.TextSize, hint)
				lb.Color = label.border.Color
				return lb.Layout(gtx)
			})
		})

	dims := layout.Stack{}.Layout(
		gtx,
		layout.Expanded(func(gtx C) D {
			// Draw a triangle to discriminate a dropdown widgets from text props.
			//      w
			//  _________  _
			//  \       /  |
			//   \  o  /   | h
			//    \   /    |
			//     \ /     |
			// (o is the offset from which we begin drawing).
			const w, h = 13, 7
			off := image.Pt(gtx.Constraints.Max.X-w, (gtx.Constraints.Min.Y)/2-h)
			stack := op.Offset(off).Push(gtx.Ops)
			anchor := clip.Path{}
			anchor.Begin(gtx.Ops)
			anchor.Move(f32.Pt(-w/2, +h/2))
			anchor.Line(f32.Pt(w, 0))
			anchor.Line(f32.Pt(-w/2, h))
			anchor.Line(f32.Pt(-w/2, -h))
			anchor.Close()
			anchorArea := clip.Outline{Path: anchor.End()}.Op()
			paint.FillShape(gtx.Ops, darkGrey, anchorArea)
			stack.Pop()

			cornerRadius := unit.Dp(4)
			dimsFunc := func(gtx C) D {
				return D{Size: image.Point{
					X: gtx.Constraints.Max.X,
					Y: gtx.Constraints.Min.Y,
				}}
			}
			border := widget.Border{
				Color:        label.border.Color,
				Width:        label.border.Thickness,
				CornerRadius: cornerRadius,
			}
			{
				visibleBorder := clip.Path{}
				visibleBorder.Begin(gtx.Ops)
				// Move from the origin to the beginning of the
				visibleBorder.LineTo(f32.Point{
					Y: float32(gtx.Constraints.Min.Y),
				})
				visibleBorder.LineTo(f32.Point{
					X: float32(gtx.Constraints.Max.X),
					Y: float32(gtx.Constraints.Min.Y),
				})
				visibleBorder.LineTo(f32.Point{
					X: float32(gtx.Constraints.Max.X),
				})
				labelStartX := float32(gtx.Dp(label.hint.Inset.Left))
				labelEndX := labelStartX + float32(label.hint.Dimensions.Size.X)
				labelEndY := float32(label.hint.Dimensions.Size.Y)
				visibleBorder.LineTo(f32.Point{
					X: labelEndX,
				})
				visibleBorder.LineTo(f32.Point{
					X: labelEndX,
					Y: labelEndY,
				})
				visibleBorder.LineTo(f32.Point{
					X: labelStartX,
					Y: labelEndY,
				})
				visibleBorder.LineTo(f32.Point{
					X: labelStartX,
				})
				visibleBorder.LineTo(f32.Point{})
				visibleBorder.Close()
				defer clip.Outline{
					Path: visibleBorder.End(),
				}.Op().Push(gtx.Ops).Pop()
			}
			return border.Layout(gtx, dimsFunc)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(unit.Dp(12)).Layout(
				gtx,
				func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return material.Label(th, th.TextSize, labelText).Layout(gtx)
				},
			)
		}),
	)

	return D{
		Size: image.Point{
			X: dims.Size.X,
			Y: dims.Size.Y + label.hint.Dimensions.Size.Y/2,
		},
		Baseline: dims.Baseline,
	}
}
