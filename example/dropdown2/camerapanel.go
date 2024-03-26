package main

import (
	l "gioui.org/layout"
	"gioui.org/unit"
	m "gioui.org/widget/material"
	c "image/color"
)

var (
	defaultTextSpacing = l.Inset{
		Top: unit.Dp(5), Bottom: unit.Dp(0),
		Left: unit.Dp(2), Right: unit.Dp(2),
	}
)

const (
	txtCameraChoice = "Вибір камери"
	txtPalette      = "Палітра"
)

type CameraPanel struct {
	th *m.Theme
}

func NewCameraPanel(theme *m.Theme) *CameraPanel {
	panel := &CameraPanel{
		th: theme,
	}
	return panel
}

func (p *CameraPanel) Layout(gtx C) D {
	ColorBox(gtx, gtx.Constraints.Max, c.NRGBA{R: 0xB0, G: 0x60, B: 0x60, A: 0xFF})

	//caption := l.Rigid(
	//	func(gtx t.C) t.D {
	//		return defaultTextSpacing.Layout(gtx, m.Overline(p.th, txtCameraChoice).Layout)
	//	},
	//)

	//var list []l.Widget

	dd := NewDropDown([]string{"w", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
	pgtx := gtx

	dimensions := l.Flex{Axis: l.Vertical}.Layout(gtx,
		l.Rigid(func(gtx l.Context) l.Dimensions {
			return dd.Layout(p.th, pgtx, gtx)
		}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
		//l.Rigid(func(gtx l.Context) l.Dimensions {
		//	return m.Label(p.th, unit.Sp(12), "a").Layout(gtx)
		//}),
	)

	//choices := func(gtx t.C) t.D {
	//	return m.ListStyle{m.List(p.th, &w.List{
	//		List: l.List{
	//			Axis:      l.Vertical,
	//			Alignment: l.Start,
	//		},
	//		Scrollbar: w.Scrollbar{},
	//	})}.Layout(gtx, len(p.rtspStreams), func(gtx t.C, i int) t.D {
	//		mp := unit.Dp(3)
	//		return l.Inset{Bottom: mp}.Layout(gtx, list[i])
	//	})
	//}

	//pageContent := append([]l.FlexChild{l.Rigid(choices)}, caption)
	//dimensions := l.Flex{Axis: l.Vertical, Alignment: l.Middle}.Layout(gtx, pageContent...)

	return dimensions
}
