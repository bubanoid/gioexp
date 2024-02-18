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
	"gioui.org/unit"
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

func NewDropDownWidget(items []string, offsetX, offsetY float64, dropdownTextSize unit.Sp, selectedMenuItemHandler MenuItemHandlerType) *DropDownWidget {
	dropdownWidget := DropDownWidget{items: items, dropdownTextSize: dropdownTextSize, selectedMenuItemHandler: selectedMenuItemHandler}
	dropdownWidget.area.OffsetX = offsetX
	dropdownWidget.area.OffsetY = offsetY
	return &dropdownWidget
}

type MenuItemHandlerType func()

type DropDownWidget struct {
	Widget
	Selected int

	items      []string
	area       ContextArea
	menu       component.MenuState
	clickables []*widget.Clickable

	focused                 bool
	click                   gesture.Click
	dropdownTextSize        unit.Sp
	selectedMenuItemHandler MenuItemHandlerType
}

func (ddWidget *DropDownWidget) Layout(th *material.Theme, pgtx, gtx C) D {
	// Handle menu selection.
	ddWidget.menu.Options = ddWidget.menu.Options[:0]
	for len(ddWidget.clickables) <= len(ddWidget.items) {
		ddWidget.clickables = append(ddWidget.clickables, &widget.Clickable{})
	}

	textSize := th.TextSize
	th.TextSize = ddWidget.dropdownTextSize
	for i := range ddWidget.items {
		click := ddWidget.clickables[i]
		if click.Clicked(gtx) {
			ddWidget.Selected = i
			ddWidget.selectedMenuItemHandler()
		}
		// todo (AA): Here we can decrease space between items in menu
		ddWidget.menu.Options = append(ddWidget.menu.Options, component.MenuItem(th, click, ddWidget.items[i]).Layout)
	}
	th.TextSize = textSize
	ddWidget.area.Activation = pointer.ButtonPrimary
	ddWidget.area.AbsolutePosition = true // todo (AA): don't clear how it works

	// Handle focus "manually". When the dropdown is closed we draw ddWidget label,
	// which can't receive focus. By registering ddWidget key.InputOp we can then receive
	// focus events (and draw the focus border). We also want to grab the focus when
	// the dropdown is opened: we do this with ddWidget.click.
	for _, e := range gtx.Events(ddWidget) {
		switch e := e.(type) {
		case key.FocusEvent:
			ddWidget.focused = e.Focus
		}
	}
	ddWidget.click.Update(gtx)

	// check if dropdown is clicked
	if ddWidget.click.Pressed() {
		// Request focus
		key.FocusOp{Tag: ddWidget}.Add(gtx.Ops)
	}

	// Clip events to the DdWidget area (area fo collapsed dropdown) only.
	clipOp := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	key.InputOp{Tag: ddWidget, Hint: key.HintAny}.Add(gtx.Ops)
	ddWidget.click.Add(gtx.Ops)
	clipOp.Pop()

	wgtx := gtx
	return layout.Stack{}.Layout(pgtx,
		layout.Stacked(func(gtx C) D {
			gtx.Constraints = layout.Exact(wgtx.Constraints.Max)

			inset := layout.Inset{Top: 1, Right: 4, Bottom: 1, Left: 4}
			label := material.Label(th, th.TextSize, ddWidget.items[ddWidget.Selected])
			label.MaxLines = 1
			label.TextSize = ddWidget.dropdownTextSize // or th.TextSize
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

			return FocusBorder(th, ddWidget.focused).Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, label.Layout)
			})
		}),
		// expanded dropdown menu
		layout.Expanded(func(gtx C) D {
			// todo (AA) th contains Inset for menu labels
			dimensions := ddWidget.area.Layout(gtx, wgtx, component.Menu(th, &ddWidget.menu).Layout)
			return dimensions
		}),
	)
}
