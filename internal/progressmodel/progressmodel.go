package progressmodel

import (
	"sync"

	icon "github.com/abaskin/Test-Maker/internal/icons"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProgressModel struct {
	modal               *widget.PopUp
	window              fyne.Window
	progressList        []string
	infiniteProgressBar bool
	progressBar         *widget.ProgressBar
	wait                sync.WaitGroup
}

type ProgressSt struct {
	Step string
	Done bool
	Err  error
	Data interface{}
}

func NewProgressModel(wind fyne.Window) *ProgressModel {
	newModel := &ProgressModel{
		progressList:        make([]string, 0),
		window:              wind,
		modal:               nil,
		infiniteProgressBar: true,
		wait:                sync.WaitGroup{},
	}
	return newModel
}

func (p *ProgressModel) SetWait() {
	p.wait.Add(1)
}

func (p *ProgressModel) Done() {
	p.wait.Done()
}

func (p *ProgressModel) Wait() {
	p.wait.Wait()
}

func (p *ProgressModel) ShowError(err error, msg ...string) {
	for _, m := range msg {
		p.Show(m, false, false)
	}
	p.Show(err.Error(), true, true)
	p.Wait()
}

func (p *ProgressModel) Show(newLine string, isError, done bool) {
	if p.modal != nil {
		p.modal.Hide()
	}
	switch isError {
	case true:
		p.showError(newLine)
	case false:
		p.progressList = append(p.progressList, newLine)
		p.modal = widget.NewModalPopUp(
			widget.NewRichTextWithText("Nothing to see here, yet..."),
			p.window.Canvas(),
		)
		p.modal.Content = p.progressBoxInfi(done)
		p.modal.Show()
	}
}

func (p *ProgressModel) Clear() *ProgressModel {
	p.progressList = make([]string, 0)
	return p
}

func (p *ProgressModel) ShowProgress(newLine string, minVal, maxVal float64) {
	if p.modal != nil {
		p.modal.Hide()
	}
	p.progressList = append(p.progressList, newLine)
	p.modal = widget.NewModalPopUp(
		widget.NewRichTextWithText("Nothing to see here, yet..."),
		p.window.Canvas(),
	)
	p.modal.Content = p.progressBox(minVal, maxVal)
	p.modal.Show()
}

func (p *ProgressModel) UpdateProgress(value float64) {
	p.progressBar.SetValue(value)
}

func (p *ProgressModel) Hide() {
	if p.modal != nil {
		p.modal.Hide()
	}
}

func (p *ProgressModel) Bar(show bool) {
	p.infiniteProgressBar = show
}

func (p *ProgressModel) showError(err string) {
	p.SetWait()
	p.modal = widget.NewModalPopUp(
		container.NewVBox(
			p.modalText("The following error has occurred."),
			p.modalText(err),
			widget.NewButton(
				"OK",
				func() {
					p.modal.Hide()
					p.wait.Done()
				},
			),
		),
		p.window.Canvas(),
	)
	p.modal.Show()
}

func (p *ProgressModel) modalText(text string) *widget.RichText {
	return widget.NewRichText(
		&widget.TextSegment{
			Text: text,
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameSubHeadingText,
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
			},
		},
	)
}

func (p *ProgressModel) progressBox(minVal, maxVal float64) *fyne.Container {
	p.progressBar = &widget.ProgressBar{
		Min:   minVal,
		Max:   maxVal,
		Value: 0,
	}
	pBox := container.NewVBox()
	for _, c := range p.makeProgressList() {
		pBox.Add(c)
	}
	pBox.Add(p.progressBar)
	return pBox
}

func (p *ProgressModel) progressBoxInfi(done bool) *fyne.Container {
	pBox := container.NewVBox()
	for _, c := range p.makeProgressList() {
		pBox.Add(c)
	}
	if p.infiniteProgressBar && !done {
		pBox.Add(widget.NewProgressBarInfinite())
	}
	if done {
		pBox.Add(
			widget.NewButton(
				"OK",
				func() {
					p.modal.Hide()
					p.wait.Done()
				},
			),
		)
	}
	return pBox
}

func (p *ProgressModel) makeProgressList() []*fyne.Container {
	pList := make([]*fyne.Container, 0)
	for _, item := range p.progressList {
		hb := container.NewHBox(
			widget.NewIcon(icon.CloudDoneOutlinedIconThemed),
			p.modalText(item),
		)
		pList = append(pList, hb)
	}
	return pList
}
