package dropdown_tunned

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
)

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

var darkGrey = rgb(0xa9a9a9)

func NewDropDownWidget(items []string) *DropDownWidget {
	return &DropDownWidget{items: items}
}

type DropDownWidget struct {
	Widget
	Selected int

	items      []string
	area       component.ContextArea
	menu       component.MenuState
	clickables []*widget.Clickable

	focused bool
	click   gesture.Click
}

func (a *DropDownWidget) Layout(th *material.Theme, pgtx, gtx C, offset image.Point) D {
	// Handle menu selection.
	a.menu.Options = a.menu.Options[:0]
	for len(a.clickables) <= len(a.items) {
		a.clickables = append(a.clickables, &widget.Clickable{})
	}
	for i := range a.items {
		click := a.clickables[i]
		if click.Clicked(gtx) {
			a.Selected = i
		}
		// todo (AA): Here we can decrease space between items in menu
		// todo (AA): We can try to change gtx here
		a.menu.Options = append(a.menu.Options, component.MenuItem(th, click, a.items[i]).Layout)
	}
	a.area.Activation = pointer.ButtonPrimary
	a.area.AbsolutePosition = true // todo (AA): don't clear how it works

	// Handle focus "manually". When the dropdown is closed we draw a label,
	// which can't receive focus. By registering a key.InputOp we can then receive
	// focus events (and draw the focus border). We also want to grab the focus when
	// the dropdown is opened: we do this with a.click.
	for _, e := range gtx.Events(a) {
		switch e := e.(type) {
		case key.FocusEvent:
			a.focused = e.Focus
		}
	}
	a.click.Update(gtx)

	// check if dropdown is clicked
	if a.click.Pressed() {
		// Request focus
		key.FocusOp{Tag: a}.Add(gtx.Ops)
		// todo (AA): simulate a click within dropdown area
		//gesture.ClickEvent{Kind: gesture.KindClick}
	}

	// Clip events to the DdWidget area only.
	// todo (AA): click within this area creates border
	clipOp := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	key.InputOp{Tag: a, Hint: key.HintAny}.Add(gtx.Ops)
	a.click.Add(gtx.Ops)
	clipOp.Pop()

	wgtx := gtx
	return layout.Stack{}.Layout(pgtx,
		layout.Stacked(func(gtx C) D {
			gtx.Constraints = layout.Exact(wgtx.Constraints.Max)
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

			inset := layout.Inset{Top: 1, Right: 4, Bottom: 1, Left: 4}
			label := material.Label(th, th.TextSize, a.items[a.Selected])
			label.MaxLines = 1
			label.TextSize = th.TextSize
			label.Alignment = text.Start
			label.Color = th.Fg

			// Draw a triangle to discriminate a dropdown widgets from text props.
			//      w
			//  _________  _
			//  \       /  |
			//   \  o  /   | h
			//    \   /    |
			//     \ /     |
			// (o is the offset from which we begin drawing).
			const w, h = 13, 7
			off := image.Pt(gtx.Constraints.Max.X-w, gtx.Constraints.Max.Y/2-h)
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

			return FocusBorder(th, a.focused).Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, label.Layout)
			})
		}),
		// expanded dropdown menu
		layout.Expanded(func(gtx C) D {
			off := op.Offset(offset).Push(gtx.Ops)
			gtx.Constraints = layout.Exact(gtx.Constraints.Max)
			// todo (AA): draw colorbox here
			// todo (AA) th contains Inset for menu labels
			dimensions := a.area.Layout(gtx, component.Menu(th, &a.menu).Layout)
			off.Pop()
			return dimensions
		}),
	)
}
