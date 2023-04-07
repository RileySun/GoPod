package main

import(
	"log"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	
	"github.com/RileySun/OmnyGo"
	
	"github.com/skratchdot/open-golang/open"
)

type List struct {
	Items []*Item
}

type Item struct {
	Name string
	FileName string
	ShortName string
	ShowName string
	AudioURL string
	Image string
	Size int64
	Downloaded bool
}

type Buttons struct {
	Download *widget.Button
	Play *widget.Button
	Delete *widget.Button
}

//Create
func NewList(slug string) *List {
	list := new(List)
	list.Items = list.getItems(slug)
	
	return list
}

func (l *List) getItems(slug string) []*Item {
	var itemList []*Item
	clips := omny.GetAllClips(slug)
	
	for _, clip := range clips {
		filename := makeFileNameSafe(clip.Title)
		item := &Item{
			Name:clip.Title,
			FileName:filename,
			ShowName:slug,
			AudioURL:clip.AudioUrl,
			Image:clip.ImageUrl,
			Size:clip.PublishedAudioSizeInBytes,
			Downloaded:checkIfDownloaded(slug + "/" + filename),
		}
		
		if len(clip.Title) > 38 {
			item.ShortName = clip.Title[0:38] + "..."
		} else {
			item.ShortName = clip.Title
		}
		
		itemList = append(itemList, item)
	}
	
	return itemList
}

//Render
func (l *List) Render() *widget.List {
	list := widget.NewList(
		func() int {
			return len(l.Items)
		},
		func() fyne.CanvasObject {
			return l.RenderItem()
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {			
			//Set Name
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(l.Items[i].ShortName)
				
			//Set Icon & Action
			buttons := Buttons{
				Download:o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button),
				Play:o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button),
				Delete:o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Button),
			}
			
			buttons.Download.OnTapped = func() {l.Download(l.Items[i], buttons)}
			buttons.Play.OnTapped = func() {l.Play(l.Items[i])}
			buttons.Delete.OnTapped = func() {l.Delete(l.Items[i], buttons)}
			
			if l.Items[i].Downloaded {
				buttons.Download.Disable()
			} else {
				buttons.Play.Disable()
				buttons.Delete.Disable()
			}
		},
	)
	
	return list
}

func (l *List) RenderItem() *fyne.Container {
	label := widget.NewLabel("Episode Name")
	
	download := widget.NewButtonWithIcon("", Icons.Download, func() {})
	play := widget.NewButtonWithIcon("", Icons.Play, func() {})
	trash := widget.NewButtonWithIcon("", Icons.Delete, func() {})
	buttons := container.New(layout.NewHBoxLayout(), download, play, trash)
	
	item := container.NewBorder(nil, nil, nil, buttons, label)
	return item
}

//Actions

func (l *List) Download(item *Item, buttons Buttons) {
	go func() {
		buttons.Download.Disable()
		downloadToFolder(item.AudioURL, item.ShowName, item.FileName)
		buttons.Play.Enable()
		buttons.Delete.Enable()
	}()
}

func (l *List) Play(item *Item) {
	path := getAppSupportFolder() + "/" + item.ShowName + "/" + item.FileName + ".mp3"
	playErr := open.Run(path)
	if playErr != nil {
		log.Fatal(playErr.Error())
	}
	
}

func (l *List) Confirm(item *Item, buttons Buttons) {
	text := "Are you sure you wish to delete\n" + item.ShortName + "?"
	dialog.ShowConfirm("Confirm Delete", text, func(b bool) {if b {l.Delete(item, buttons)}}, window)
}

func (l *List) Delete(item *Item, buttons Buttons) {
	buttons.Play.Disable()
	buttons.Delete.Disable()
	deleteFile(item.ShowName, item.FileName)
	buttons.Download.Enable()
}