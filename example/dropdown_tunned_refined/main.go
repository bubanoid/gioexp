package main

import (
	"gioui.org/app"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/constraints"
	"image"
	"log"
	"os"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	DefaultPropertyHeight     = unit.Dp(30)
	DefaultPropertyWidth      = unit.Dp(140)
	DefaultPropertyListHeight = unit.Dp(100)
	DefaultOffsetX            = 200
	DefaultOffsetY            = 0
)

type UI struct {
	th *material.Theme
	dd *DropDown
}

// A DropDown holds and presents a vertical, scrollable list of properties. A DropDown
// is divided into 2 columns: property names on the left and widgets for
// property values on the right. These 2 sections can be resized thanks to a
// divider, which can be dragged.
type DropDown struct {
	DdWidget DropDownWidget

	// PropertyListHeight is the height of a expanded list
	PropertyListHeight unit.Dp

	// PropertyHeight is the height of a property. All properties have
	// the same dimensions. The width depends on the horizontal space available
	// for the list
	PropertyHeight unit.Dp

	// PropertyWidth is the width of a property. All properties have
	// the same width.
	PropertyWidth unit.Dp

	// offset is the offset of the dropdown values
	Offset image.Point
}

// NewDropdown creates a new DropDown.
func NewDropdown(ddValues []string) *DropDown {
	return &DropDown{
		DdWidget:           *NewDropDownWidget(ddValues),
		PropertyListHeight: DefaultPropertyListHeight,
		PropertyHeight:     DefaultPropertyHeight,
		PropertyWidth:      DefaultPropertyWidth,
		Offset:             image.Point{X: DefaultOffsetX, Y: DefaultOffsetY},
	}
}

func (dd *DropDown) visibleListHeight(gtx C) int {
	return min(gtx.Dp(dd.PropertyListHeight), gtx.Constraints.Max.Y)
}

func (dd *DropDown) visibleHeight(gtx C) int {
	return min(gtx.Dp(dd.PropertyHeight), gtx.Constraints.Max.Y)
}

func (dd *DropDown) visibleWidth(gtx C) int {
	return min(gtx.Dp(dd.PropertyWidth), gtx.Constraints.Max.X)
}

func (dd *DropDown) Layout(th *material.Theme, gtx C) D {
	wtotal := dd.visibleWidth(gtx)
	htotal := dd.visibleHeight(gtx)
	hlist := dd.visibleListHeight(gtx)

	gtx.Constraints.Max.X = wtotal

	dim := widget.Border{
		Color:        th.Fg,
		CornerRadius: unit.Dp(2),
		Width:        unit.Dp(1),
	}.Layout(gtx, func(gtx C) D {
		// Copy the context passed to property widgets, we don't want
		// its size constrained since it's used as modal pane.
		pgtx := gtx
		pgtx.Constraints = layout.Exact(image.Pt(wtotal, htotal+hlist))
		gtx.Constraints = layout.Exact(image.Pt(wtotal, htotal))

		return layout.Inset{}.Layout(
			gtx,
			func(gtx C) D {
				gtx.Constraints.Min.Y = gtx.Dp(dd.PropertyHeight)
				gtx.Constraints.Max.Y = gtx.Dp(dd.PropertyHeight)

				// Draw dropdown
				size := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
				gtx.Constraints = layout.Exact(size)
				dd.DdWidget.Layout(th, pgtx, gtx, dd.Offset)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			},
		)
	})

	return dim
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Widget shows the value of a property and handles user actions to edit it.
type Widget interface {
	// Layout lays out the property DdWidget using gtx which is the
	// property-specific context, and pgtx which is the parent context (useful
	// for properties that require more space during edition).
	Layout(th *material.Theme, pgtx, gtx layout.Context) D
}

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
	}

	// Clip events to the DdWidget area only.
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

			return inset.Layout(gtx, label.Layout)
		}),
		// expanded dropdown menu
		layout.Expanded(func(gtx C) D {
			off := op.Offset(offset).Push(gtx.Ops)
			gtx.Constraints = layout.Exact(gtx.Constraints.Max)
			dimensions := a.area.Layout(gtx, component.Menu(th, &a.menu).Layout)
			off.Pop()
			return dimensions
		}),
	)
}

func NewUI(theme *material.Theme) *UI {
	ui := &UI{
		th: theme,
	}

	dropdown := NewDropdown(
		[]string{"It", "is", "not", "easy", "to", "work", "with", "gio"},
	)

	ui.dd = dropdown
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	return ui.dd.Layout(ui.th, gtx)
}

func main() {
	ui := NewUI(material.NewTheme())

	go func() {
		w := app.NewWindow(app.Title("DropDown Menu"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
