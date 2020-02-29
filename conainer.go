package tview

import (
	"github.com/gdamore/tcell"
)

// Configuration values.
const (
	ContainerRow = iota
	ContainerColumn
)

// ContainerItem holds layout options for one item.
type containerItem struct {
	Item Primitive // The item to be positioned. May be nil for an empty item.

	// Fixed size. If < 0, then the primitive size will be changed
	// according to its needs
	FixedSize int
	Focus     bool // Whether or not this item attracts the layout's focus.
}

// Container allows to stack primitives one after another. It uses their required size
// or fixed set by AddItem
type Container struct {
	*Box

	// The items to be positioned.
	items []*containerItem

	// ContainerRow or ContainerColumn.
	direction int

	// If set to true, Container will use the entire screen as its available space
	// instead its box dimensions.
	fullScreen bool
}

// NewContainer returns a new Container with no primitives and its
// direction set to ContainerColumn. To add primitives to this layout, see AddItem().
// To change the direction, see SetDirection().
//
// Note that Box, the superclass of Container, will have its background color set to
// transparent so that any nil Container items will leave their background unchanged.
// To clear a Container's background before any items are drawn, set it to the
// desired color:
//
//   Container.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
func NewContainer() *Container {
	f := &Container{
		Box:       NewBox().SetBackgroundColor(tcell.ColorDefault),
		direction: ContainerColumn,
	}
	f.focus = f
	return f
}

// SetDirection sets the direction in which the contained primitives are
// distributed. This can be either ContainerColumn (default) or ContainerRow.
func (f *Container) SetDirection(direction int) *Container {
	f.direction = direction
	return f
}

// SetFullScreen sets the flag which, when true, causes the Container layout to use
// the entire screen space instead of whatever size it is currently assigned to.
func (f *Container) SetFullScreen(fullScreen bool) *Container {
	f.fullScreen = fullScreen
	return f
}

// AddItem adds a new item to the container. If fixedSize < 0, the item will be resized
// acording to its needs (by calling GetWidth or GetHeight).
// If "focus" is set to true, the item will receive focus when the Container
// primitive receives focus. If multiple items have the "focus" flag set to
// true, the first one will receive focus.
//
// You can provide a nil value for the primitive. This will still consume screen
// space but nothing will be drawn.
func (f *Container) AddItem(item Primitive, fixedSize int, focus bool) *Container {
	if item != nil {
		f.items = append(f.items, &containerItem{Item: item, FixedSize: fixedSize, Focus: focus})
	}
	return f
}

// RemoveItem removes all items for the given primitive from the container,
// keeping the order of the remaining items intact.
func (f *Container) RemoveItem(p Primitive) *Container {
	for index := len(f.items) - 1; index >= 0; index-- {
		if f.items[index].Item == p {
			f.items = append(f.items[:index], f.items[index+1:]...)
		}
	}
	return f
}

// Clear removes all items from the container.
func (f *Container) Clear() *Container {
	f.items = nil
	return f
}

// ResizeItem sets a new size for the item(s) with the given primitive. If there
// are multiple Container items with the same primitive, they will all receive the
// same size. For details regarding the size parameters, see AddItem().
func (f *Container) ResizeItem(p Primitive, fixedSize, proportion int) *Container {
	for _, item := range f.items {
		if item.Item == p {
			item.FixedSize = fixedSize
		}
	}
	return f
}

// Draw draws this primitive onto the screen.
func (f *Container) Draw(screen tcell.Screen) {
	f.Box.Draw(screen)

	// Calculate size and position of the items.

	// Do we use the entire screen?
	if f.fullScreen {
		width, height := screen.Size()
		f.SetRect(0, 0, width, height)
	}

	// How much space can we distribute?
	x, y, width, height := f.GetInnerRect()

	// Calculate positions and draw items.
	pos := x
	maxPos := x + width
	if f.direction == ContainerRow {
		pos = y
		maxPos = y + height
	}
	for _, item := range f.items {
		// Draw empty item if its has fixed size
		if item.Item == nil {
			if item.FixedSize >= 0 {
				pos += item.FixedSize
			}
			continue
		}
		// Get the required size
		var size int
		if item.FixedSize >= 0 {
			size = item.FixedSize
		} else if f.direction == ContainerRow {
			size = item.Item.GetHeight(width)
		} else {
			size = item.Item.GetWidth(height)
		}
		if pos+size > maxPos {
			size = maxPos - pos
		}
		// Draw the primitive
		if size > 0 {
			if f.direction == ContainerColumn {
				item.Item.SetRect(pos, y, size, height)
			} else {
				item.Item.SetRect(x, pos, width, size)
			}

			pos += size
			if item.Item.GetFocusable().HasFocus() {
				defer item.Item.Draw(screen)
			} else {
				item.Item.Draw(screen)
			}
		}
		if pos > maxPos {
			break
		}
	}
}

// Focus is called when this primitive receives focus.
func (f *Container) Focus(delegate func(p Primitive)) {
	for _, item := range f.items {
		if item.Item != nil && item.Focus {
			delegate(item.Item)
			return
		}
	}
}

// HasFocus returns whether or not this primitive has focus.
func (f *Container) HasFocus() bool {
	for _, item := range f.items {
		if item.Item != nil && item.Item.GetFocusable().HasFocus() {
			return true
		}
	}
	return false
}
