package main

import(
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var fyneApp fyne.App
var window fyne.Window

var list *List

func setup() {
	fyneApp = app.NewWithID("com.sunshine.gopod")
	window = fyneApp.NewWindow("Go Pod")
	
	
	list = NewList("cool-people-who-did-cool-stuff")
	
	/*
	//System Tray Notifications (desktop only, no fyne tray for mobile)
	device := fyne.CurrentDevice()
	if !device.IsMobile() {
		systemTray = //NEW TRAY CODE HERE
	}
	*/
}

func main() {
	log.Println("Loading Go Pod")
	setup()
	log.Println("Loading Complete")
	
	content := list.Render()
	
	window.SetContent(content)
	
	window.CenterOnScreen()
	
	window.Resize(fyne.NewSize(400, 600))
	window.SetFixedSize(true)
	
	window.SetMaster()
	
	window.ShowAndRun()
}