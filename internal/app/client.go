package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

func RunClient() {
	app := tview.NewApplication()

	textArea := tview.NewTextArea().
		SetPlaceholder("Enter msg here...")
	textArea.SetTitle("Input").SetBorder(true)

	chat := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).ScrollToEnd()
	chat.SetBorder(true).SetTitle("Chat")

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

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Users"), 0, 3, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(chat, 0, 10, false).
			AddItem(textArea, 0, 2, false), 0, 10, false)
	if err := app.SetRoot(flex, true).SetFocus(textArea).Run(); err != nil {
		panic(err)
	}
}
