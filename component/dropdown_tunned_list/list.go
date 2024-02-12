package dropdown_tunned_list

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
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
	DefaultHandleBarWidth = unit.Dp(3)
	DefaultPropertyWidth  = unit.Dp(120)
)

// A List holds and presents a vertical, scrollable list of properties. A List
// is divided into 2 columns: property names on the left and widgets for
// property values on the right. These 2 sections can be resized thanks to a
// divider, which can be dragged.
type List struct {
	widgets []Widget

	// PropertyHeight is the height of a single property. All properties have
	// the same dimensions. The width depends on the horizontal space available
	// for the list
	PropertyHeight unit.Dp

	PropertyWidth unit.Dp

	// HandleBarWidth is the width of the handlebar used to resize the columns.
	HandleBarWidth unit.Dp

	list layout.List

	// ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	//ratio float32
}

// NewList creates a new List.
func NewList() *List {
	return &List{
		PropertyHeight: DefaultPropertyHeight,
		PropertyWidth:  DefaultPropertyWidth,
		HandleBarWidth: DefaultHandleBarWidth,
		list: layout.List{
			Axis: layout.Vertical,
		},
		//ratio: -1.0,
	}
}

// Add adds a new property to the list.
func (plist *List) Add(widget Widget) {
	plist.widgets = append(plist.widgets, widget)
}

func (plist *List) visibleHeight(gtx C) int {
	return min(gtx.Dp(plist.PropertyHeight)*len(plist.widgets), gtx.Constraints.Max.Y)
}

func (plist *List) visibleWidth(gtx C) int {
	return min(gtx.Dp(plist.PropertyWidth)*len(plist.widgets), gtx.Constraints.Max.X)
}

func (plist *List) Layout(th *material.Theme, gtx C) D {
	htotal := plist.visibleHeight(gtx)
	wtotal := plist.visibleWidth(gtx)

	dim := widget.Border{
		Color:        th.Fg,
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}.Layout(gtx, func(gtx C) D {
		pgtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(wtotal, htotal))
		return plist.list.Layout(gtx, len(plist.widgets), func(gtx C, i int) D {
			gtx.Constraints.Min.Y = gtx.Dp(plist.PropertyHeight)
			gtx.Constraints.Max.Y = gtx.Dp(plist.PropertyHeight)

			rsize := gtx.Constraints.Max.X

			// Draw dropdown value.
			off := op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
			size := image.Pt(rsize, gtx.Constraints.Max.Y)
			gtx.Constraints = layout.Exact(size)
			plist.widgets[i].Layout(th, pgtx, gtx)
			off.Pop()

			return layout.Dimensions{Size: gtx.Constraints.Max}
		})
	})

	return dim
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func clamp[T constraints.Ordered](mn, val, mx T) T {
	if val < mn {
		return mn
	}
	if val > mx {
		return mx
	}
	return val
}

// Widget shows the value of a property and handles user actions to edit it.
type Widget interface {
	// Layout lays out the property widget using gtx which is the
	// property-specific context, and pgtx which is the parent context (useful
	// for properties that require more space during edition).
	Layout(th *material.Theme, pgtx, gtx layout.Context) D
}
