package main

import(
	"fyne.io/fyne/v2/widget"
)

type Show struct {
	Type string
	Name string
	Slug string
}

type Menu struct {
	Shows []*Show
	Select *widget.Select
}

func NewMenu(shows []*Show, changeCallback func(string)) *Menu {
	var showStrings []string
	for _, show := range shows {
		showStrings = append(showStrings, show.Slug)
	}
	
	menu := &Menu {
		Shows:shows,
		Select:widget.NewSelect(showStrings, changeCallback),
	}
	
	return menu
}
