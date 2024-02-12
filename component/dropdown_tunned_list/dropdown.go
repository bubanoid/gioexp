package dropdown_tunned_list

import (
	"image"

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
)

var darkGrey = rgb(0xa9a9a9)

func NewDropDown(items []string) *DropDown {
	return &DropDown{items: items}
}

type DropDown struct {
	Selected int

	items      []string
	area       component.ContextArea
	menu       component.MenuState
	clickables []*widget.Clickable

	focused bool
	click   gesture.Click
}

func (dd *DropDown) Layout(th *material.Theme, pgtx, gtx C) D {
	// Handle menu selection.
	dd.menu.Options = dd.menu.Options[:0]
	for len(dd.clickables) <= len(dd.items) {
		dd.clickables = append(dd.clickables, &widget.Clickable{})
	}
	for i := range dd.items {
		click := dd.clickables[i]
		if click.Clicked(gtx) {
			dd.Selected = i
		}
		dd.menu.Options = append(dd.menu.Options, component.MenuItem(th, click, dd.items[i]).Layout)
	}
	dd.area.Activation = pointer.ButtonPrimary
	dd.area.AbsolutePosition = true

	// Handle focus "manually". When the dropdown is closed we draw dd label,
	// which can't receive focus. By registering dd key.InputOp we can then receive
	// focus events (and draw the focus border). We also want to grab the focus when
	// the dropdown is opened: we do this with dd.click.
	for _, e := range gtx.Events(dd) {
		switch e := e.(type) {
		case key.FocusEvent:
			dd.focused = e.Focus
		}
	}
	dd.click.Update(gtx)
	if dd.click.Pressed() {
		// Request focus
		key.FocusOp{Tag: dd}.Add(gtx.Ops)
	}

	// Clip events to the widget area only.
	clipOp := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	key.InputOp{Tag: dd, Hint: key.HintAny}.Add(gtx.Ops)
	dd.click.Add(gtx.Ops)
	clipOp.Pop()

	wgtx := gtx
	return layout.Stack{}.Layout(pgtx,
		layout.Stacked(func(gtx C) D {
			gtx.Constraints = layout.Exact(wgtx.Constraints.Max)
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

			inset := layout.Inset{Top: 1, Right: 4, Bottom: 1, Left: 4}
			label := material.Label(th, th.TextSize, dd.items[dd.Selected])
			label.MaxLines = 1
			label.TextSize = th.TextSize
			label.Alignment = text.Start
			label.Color = th.Fg

			// Draw dd triangle to discriminate dd drop down widgets from text props.
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

			return FocusBorder(th, dd.focused).Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, label.Layout)
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints = layout.Exact(gtx.Constraints.Max)
			return dd.area.Layout(gtx, component.Menu(th, &dd.menu).Layout)
		}),
	)
}
