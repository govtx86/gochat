// Package app provides the main application.
package app

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var conn net.Conn = nil

func RunClient() {
	var err error
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
	chat.SetChangedFunc(func() {
		app.Draw()
	})

	users := tview.NewTextView().SetDynamicColors(true)
	users.SetBorder(true).SetTitle("Users")
	users.SetChangedFunc(func() {
		app.Draw()
	})

	flex := tview.NewFlex().
		AddItem(users, 0, 3, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(chat, 0, 10, false).
			AddItem(textArea, 0, 2, true), 0, 10, false)

	connection := tview.NewForm()

	pages := tview.NewPages()
	pages.AddPage("connect", modal(connection, 50, 13), true, true)
	pages.AddPage("chat", flex, true, false)

	connection.
		AddTextView("Status", "Waiting", 0, 1, true, false).
		AddInputField("Username", "", 25, nil, nil).
		AddInputField("Address", "", 25, nil, nil).
		AddInputField("Port", "8080", 10, nil, nil)

	connection.
		AddButton("Connect", func() {
			username := connection.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			address := connection.GetFormItemByLabel("Address").(*tview.InputField).GetText()
			port := connection.GetFormItemByLabel("Port").(*tview.InputField).GetText()
			conn, err = net.Dial("tcp", net.JoinHostPort(address, port))
			if err != nil {
				connection.GetFormItemByLabel("Status").(*tview.TextView).SetText("Error Connecting")
				return
			}
			conn.Write([]byte(username + "\n"))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			status := strings.Trim(response, "\r\n")
			switch status {
			case "409":
				connection.GetFormItemByLabel("Status").(*tview.TextView).SetText("Username already in use")
			case "200":
				connection.GetFormItemByLabel("Status").(*tview.TextView).SetText("Connected")
				pages.SwitchToPage("chat")
				go func() {
					for {
						msgResp, err := bufio.NewReader(conn).ReadString('\n')
						if err != nil {
							app.Stop()
							return
						}
						if strings.Contains(msgResp, "#srvc:") {
							go func() {
								app.QueueUpdateDraw(func() {
									users.SetText(strings.ReplaceAll(strings.Trim(msgResp, "#srvc:"), "#$", "\n"))
								})
							}()
						} else {
							msg := strings.Trim(msgResp, "\r\n")
							chat.SetText(chat.GetText(true) + msg + "\n")
						}
					}
				}()
			default:
				connection.GetFormItemByLabel("Status").(*tview.TextView).SetText("Err0r Connecting")
			}
		}).
		AddButton("Cancel", func() {
			app.Stop()
		}).SetButtonsAlign(1)

	connection.SetBorder(true).SetTitle("Network Connection")
	// TextArea handle enter
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			msg := strings.TrimSpace(textArea.GetText())
			if msg != "" {
				conn.Write([]byte(msg + "\n"))
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
		fmt.Println("Error running app:", err)
	}
}
