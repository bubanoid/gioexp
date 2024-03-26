package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	l "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	m "gioui.org/widget/material"
	"github.com/arl/gioexp/component/qasset"
	"image"
	c "image/color"
)

type (
	C = l.Context
	D = l.Dimensions
)

var (
	Black        = c.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	MidnightBlue = c.NRGBA{R: 0x19, G: 0x19, B: 0x70, A: 0xFF}
	White        = c.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	Indigo       = c.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
	AliceBlue    = c.NRGBA{R: 240, G: 248, B: 255, A: 255}
	DarkGrey     = rgb(0xa9a9a9) // todo (AA) use RGB
	LightGrey    = rgb(0xd3d3d3)
)

func rgb(c uint32) c.NRGBA {
	return argb(0xff000000 | c)
}

func argb(cuint uint32) c.NRGBA {
	return c.NRGBA{A: uint8(cuint >> 24), R: uint8(cuint >> 16), G: uint8(cuint >> 8), B: uint8(cuint)}
}

type Panel struct {
	th *m.Theme
}

func NewP(theme *m.Theme) *Panel {
	panel := &Panel{th: theme}
	return panel
}

func (p *Panel) Layout(_ l.Context, _ *m.Theme) l.Dimensions {
	return l.Dimensions{}
}

type UI struct {
	w     *app.Window
	panel *Panel
	theme *m.Theme
	cc    *CameraPanel
}

func New(w *app.Window) *UI {
	theme := fixThemePalette(m.NewTheme())
	u := &UI{
		w:     w,
		theme: theme,
		panel: NewP(theme),
		cc:    NewCameraPanel(theme),
	}
	return u
}

func ColorBox(gtx l.Context, size image.Point, color c.NRGBA) l.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return l.Dimensions{Size: size}
}

func vimgLayout(gtx C) D {
	marginTopBottom, marginLeftRight := 800, 600

	return l.Flex{Axis: l.Vertical}.Layout(gtx,
		l.Flexed(1, func(gtx C) D {
			return D{Size: image.Point{Y: marginTopBottom}}
		}),
		l.Rigid(func(gtx C) D {
			return l.Flex{Axis: l.Horizontal}.Layout(gtx,
				l.Flexed(1, func(gtx C) D {
					return D{Size: image.Point{X: marginLeftRight}}
				}),
				l.Rigid(func(gtx C) D {
					imageOp := paint.NewImageOp(qasset.Neutral)
					return widget.Image{
						Src:      imageOp,
						Fit:      widget.Cover,
						Position: l.Center,
					}.Layout(gtx)
				}),
			)
		}),
	)
}

func (u *UI) Layout(gtx l.Context) l.Dimensions {
	ColorBox(gtx, gtx.Constraints.Max, Black)
	dimensions := l.Flex{Axis: l.Vertical, Spacing: l.SpaceEnd}.Layout(gtx,
		// todo (AA): here
		l.Flexed(1, vimgLayout),
		l.Rigid(
			func(gtx l.Context) l.Dimensions {
				ColorBox(gtx, gtx.Constraints.Max, c.NRGBA{R: 0x0, G: 0x0, B: 0xC0, A: 0xFF})
				return l.Flex{Axis: l.Horizontal, Alignment: l.Start, Spacing: l.SpaceBetween}.Layout(gtx,
					l.Rigid(
						func(gtx l.Context) l.Dimensions {
							return l.Flex{Axis: l.Horizontal, Alignment: l.Start, Spacing: l.SpaceStart}.Layout(gtx,
								l.Rigid(
									func(gtx l.Context) l.Dimensions {
										return u.cc.Layout(gtx)
									},
								),
							)
						},
					),
				)
			},
		),
	)
	return dimensions
}

func fixThemePalette(src *m.Theme) *m.Theme {
	if src == nil {
		return nil
	}
	palette := m.Palette{
		Bg:         c.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Fg:         c.NRGBA{R: 255, G: 255, B: 255, A: 255},
		ContrastBg: c.NRGBA{R: 63, G: 81, B: 181, A: 255},
		ContrastFg: c.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}
	dst := src.WithPalette(palette)
	dst.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	dst.FingerSize = 30

	return &dst
}
