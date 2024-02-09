package dropdown_tunned

import (
	"gioui.org/op"
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/constraints"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	DefaultPropertyHeight = unit.Dp(30)
	DefaultPropertyWidth  = unit.Dp(160)
)

// A DropDown holds and presents a vertical, scrollable list of properties. A DropDown
// is divided into 2 columns: property names on the left and widgets for
// property values on the right. These 2 sections can be resized thanks to a
// divider, which can be dragged.
type DropDown struct {
	DdWidget DropDownWidget

	// PropertyHeight is the height of a property. All properties have
	// the same dimensions. The width depends on the horizontal space available
	// for the list
	PropertyHeight unit.Dp

	// PropertyWidth is the width of a property. All properties have
	// the same width.
	PropertyWidth unit.Dp

	// offset is the offset of the dropdown values
	offset image.Point
}

// NewDropdown creates a new DropDown.
func NewDropdown(ddWidget []string) *DropDown {
	return &DropDown{
		DdWidget:       *NewDropDownWidget(ddWidget),
		PropertyHeight: DefaultPropertyHeight,
		PropertyWidth:  DefaultPropertyWidth,
	}
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

	gtx.Constraints.Max.X = wtotal

	dim := widget.Border{
		Color:        th.Fg,
		CornerRadius: unit.Dp(2),
		Width:        unit.Dp(1),
	}.Layout(gtx, func(gtx C) D {
		// Copy the context passed to property widgets, we don't want
		// its size constrained since it's used as modal pane.
		pgtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(wtotal, htotal))

		return layout.Inset{}.Layout(
			gtx,
			func(gtx C) D {
				gtx.Constraints.Min.Y = gtx.Dp(dd.PropertyHeight)
				gtx.Constraints.Max.Y = gtx.Dp(dd.PropertyHeight)
				return dd.layoutProperty(th, pgtx, gtx)
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

// layoutProperty lays out the property at index i from the list.
func (dd *DropDown) layoutProperty(th *material.Theme, pgtx, gtx C) D {
	rsize := gtx.Constraints.Max.X

	{
		// Draw property value.
		gtx := gtx
		off := op.Offset(dd.offset).Push(gtx.Ops)
		size := image.Pt(rsize, gtx.Constraints.Max.Y)
		gtx.Constraints = layout.Exact(size)
		dd.DdWidget.Layout(th, pgtx, gtx)
		off.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// Widget shows the value of a property and handles user actions to edit it.
type Widget interface {
	// Layout lays out the property DdWidget using gtx which is the
	// property-specific context, and pgtx which is the parent context (useful
	// for properties that require more space during edition).
	Layout(th *material.Theme, pgtx, gtx layout.Context) D
}
