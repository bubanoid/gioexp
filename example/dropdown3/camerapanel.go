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

	dd := dropdown_tunned_fail.NewDropdown(
		[]string{"ciao", "bonjour", "hello", "hallo", "buongiorno", "buenos dias", "ola", "bom dia"},
	)

	dimensions := l.Flex{Axis: l.Vertical}.Layout(gtx,
		l.Rigid(func(gtx l.Context) l.Dimensions {
			return dd.Layout(p.th, gtx)
		}),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "a123").Layout(gtx) }),
		l.Rigid(func(gtx l.Context) l.Dimensions { return m.Label(p.th, unit.Sp(12), "b123").Layout(gtx) }),
	)

	return dimensions
}
