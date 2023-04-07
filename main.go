package main

import(
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var fyneApp fyne.App
var window fyne.Window
var content *fyne.Container

var menu *Menu
var list *List

func setup() {
	fyneApp = app.NewWithID("com.sunshine.gopod")
	window = fyneApp.NewWindow("Go Pod")
	
	shows := loadShowsFromJSON()
	menu = NewMenu(shows, changeList)
	list = NewList(shows[0].Slug)
	
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
	
	
	content := container.NewBorder(menu.Select, nil, nil, nil, list.Render())
	
	window.SetContent(content)
	
	window.CenterOnScreen()
	
	window.Resize(fyne.NewSize(400, 600))
	window.SetFixedSize(true)
	
	window.SetMaster()
	
	window.ShowAndRun()
}

func changeList(newShow string) {
	list.Items = list.getItems(newShow)
	content := container.NewBorder(menu.Select, nil, nil, nil, list.Render()) 
	window.SetContent(content)
}