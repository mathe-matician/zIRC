package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	inputField := tview.NewInputField().
		SetLabel("Enter a number: ").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		}).
		SetBorder(true).  // Set border directly on the InputField
		SetTitle("Input") // Set the title for the InputField

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Channels/Convos"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Chat window"), 0, 3, false).
			AddItem(inputField.SetBackgroundColor(tcell.ColorLightPink), 5, 1, true), 0, 4, false)

	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
}
