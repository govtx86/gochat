// Package app provides the main application.
package app

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunClient() {
	app := tview.NewApplication().EnableMouse(true)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	textArea := tview.NewTextArea().
		SetPlaceholder("Enter msg here...")
	textArea.SetTitle("Input").SetBorder(true)

	chat := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).ScrollToEnd()
	chat.SetBorder(true).SetTitle("Chat")

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Users"), 0, 3, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(chat, 0, 10, false).
			AddItem(textArea, 0, 2, true), 0, 10, false)

	connection := tview.NewForm()

	connection.
		AddTextView("Status", "Waiting", 0, 1, true, false).
		AddInputField("Address", "", 25, nil, nil).
		AddInputField("Port", "", 10, nil, nil).
		AddButton("Connect", func() {
			// Connect to server
			connection.GetFormItemByLabel("Status").(*tview.TextView).SetText("Connected")
		}).
		AddButton("Cancel", func() {
			app.Stop()
		}).SetButtonsAlign(1)

	connection.SetBorder(true).SetTitle("Network Connection")

	pages := tview.NewPages()
	pages.AddPage("connect", modal(connection, 50, 12), true, true)
	pages.AddPage("chat", flex, true, false)

	// TextArea handle enter
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			if strings.TrimSpace(textArea.GetText()) != "" {
				chat.SetText(chat.GetText(true) + strings.TrimSpace(textArea.GetText()) + "\n")
			}
			go func() {
				app.QueueUpdateDraw(func() {
					textArea.SetText("", false)
				})
			}()
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(textArea).Run(); err != nil {
		panic(err)
	}
}
