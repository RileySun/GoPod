package main

import(
	"log"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	
	"github.com/RileySun/OmnyGo"
)

type List struct {
	Items []*Item
}

type Item struct {
	Name string
	FileName string
	ShortName string
	AudioURL string
	Image string
	Size int64
	Downloaded bool
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
			AudioURL:clip.AudioUrl,
			Image:clip.ImageUrl,
			Size:clip.PublishedAudioSizeInBytes,
			Downloaded:checkIfDownloaded(filename),
		}
		
		if len(clip.Title) > 50 {
			item.ShortName = clip.Title[0:50] + "..."
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
			button := o.(*fyne.Container).Objects[1].(*widget.Button)
			if l.Items[i].Downloaded {
				button.SetIcon(Icons.Delete)
				button.OnTapped = func() {l.Confirm(l.Items[i])}
			} else {
				button.OnTapped =  func() {l.Download(l.Items[i])}
			}
		},
	)
	
	return list
}

func (l *List) RenderOLD() *widget.List {
	list := widget.NewList(
		func() int {
			return len(l.Items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//Trunc if too long (and add ellipses)
			var name string
			rawName := l.Items[i].Name
			if len(rawName) > 55 {
				name = rawName[0:55] + "..."
			} else {
				name = rawName
			}
			_ = name
			
			//o.(*widget.Label).SetText(name)
		},
	)
	
	return list
}

func (l *List) RenderItem() *fyne.Container {
	label := widget.NewLabel("Episode Name")
	button := widget.NewButtonWithIcon("", Icons.Download, func() {})
	item := container.NewBorder(nil, nil, nil, button, label)
	return item
}

func (l *List) Download(item *Item) {
	//log.Println("Download")
	log.Println(item.FileName)
}

func (l *List) Confirm(item *Item) {
	text := "Are you sure you wish to delete\n" + item.ShortName + "?"
	dialog.ShowConfirm("Confirm Delete", text, func(b bool) {if b {l.Delete(item)}}, window)
}

func (l *List) Delete(item *Item) {
	log.Println("Delete")
}