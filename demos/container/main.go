// Demo code for the Container primitive.
package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Usage:

// up/down arrows to select textview
// type to append to the current textview
// press delete to clear current textview

func main() {
	var textViews []*tview.TextView
	current := 0
	app := tview.NewApplication()
	container := tview.NewContainer()

	flex := tview.NewFlex()
	tv := tview.NewTextView()
	tv.SetBorder(true).SetBorderColor(tcell.ColorRed)
	textViews = append(textViews, tv)
	box := tview.NewBox().SetBackgroundColor(tcell.ColorGreen)
	flex.AddItem(box, 1, 0, false)
	flex.AddItem(tv, -1, 1, false)

	flexVert := tview.NewFlex().SetDirection(tview.FlexRow)
	tv = tview.NewTextView()
	textViews = append(textViews, tv)
	tv.SetBorder(true)

	flexVert.AddItem(nil, -1, 1, false)
	flexVert.AddItem(tv, -1, 1, false)
	flexVert.AddItem(nil, -1, 1, false)
	flex.AddItem(flexVert, -1, 1, false)
	container.AddItem(flex, -1, false).SetDirection(tview.ContainerRow).SetBorder(true)
	container.SetRect(0, 0, 80, 35)

	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})

	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Rune() == '+':
			tv := tview.NewTextView()
			tv.SetBorder(true)
			container.AddItem(tv, -1, false)
			textViews = append(textViews, tv)
		case event.Key() == tcell.KeyUp:
			textViews[current].SetBorderColor(tcell.ColorDefault)
			current--
			if current < 0 {
				current = 0
			}
			textViews[current].SetBorderColor(tcell.ColorRed)
		case event.Key() == tcell.KeyDown:
			textViews[current].SetBorderColor(tcell.ColorDefault)
			current++
			if current >= len(textViews) {
				current = len(textViews) - 1
			}
			textViews[current].SetBorderColor(tcell.ColorRed)
		case event.Key() == tcell.KeyDelete:
			textViews[current].Clear()
		case event.Key() == tcell.KeyEnter:
			fmt.Fprintln(textViews[current])
		case event.Rune() != 0:
			fmt.Fprintf(textViews[current], "%c", event.Rune())
		}
		return nil
	})

	if err := app.SetRoot(container, false).Run(); err != nil {
		panic(err)
	}
}
