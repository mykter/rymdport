package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

var emptyItem = &SendItem{}

// SendItem is the item that is being sent.
type SendItem struct {
	Progress *SendProgress
	Code     chan string
	URI      fyne.URI
}

// ProgressList is a list of progress bars that track send progress.
type ProgressList struct {
	widget.List

	Items []SendItem
}

// Length returns the length of the data.
func (p *ProgressList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *ProgressList) CreateItem() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(newSendLayout(), widget.NewFileIcon(nil), widget.NewLabel("Waiting for filename..."), newCodeDisplay(), NewSendProgress())
}

// UpdateItem updates the data in the list.
func (p *ProgressList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].URI.Name())
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*CodeDisplay).waitForCode(p.Items[i].Code)
	p.Items[i].Progress = item.(*fyne.Container).Objects[3].(*SendProgress)
}

// OnItemSelected handles removing items and stopping send (in the future)
func (p *ProgressList) OnItemSelected(i int) {
	if p.Items[i].Progress.Value != p.Items[i].Progress.Max { // TODO: Stop the send instead.
		return // We can't stop running sends due to bug in wormhole-gui.
	}

	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			// Make sure that GC run on removed element
			copy(p.Items[i:], p.Items[i+1:])
			p.Items[p.Length()-1] = *emptyItem
			p.Items = p.Items[:p.Length()-1]

			p.Refresh()
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// NewSendItem adds data about a new send to the list and then returns the channel to update the code.
func (p *ProgressList) NewSendItem(URI fyne.URI) chan string {
	p.Items = append(p.Items, SendItem{Progress: NewSendProgress(), URI: URI, Code: make(chan string)})
	p.Refresh()

	return p.Items[p.Length()-1].Code
}

// NewProgressList greates a list of progress bars.
func NewProgressList() *ProgressList {
	p := &ProgressList{}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnItemSelected = p.OnItemSelected
	p.ExtendBaseWidget(p)

	return p
}
