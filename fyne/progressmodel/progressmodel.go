package progressmodel

import (
	"sync"

	icon "github.com/abaskin/testparts/fyne/icons"

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
	waitCount           uint8
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
		waitCount:           0,
	}
	return newModel
}

func (p *ProgressModel) SetWait() *ProgressModel {
	p.wait.Add(1)
	p.waitCount += 1
	return p
}

func (p *ProgressModel) Done() *ProgressModel {
	p.wait.Done()
	p.waitCount -= 1
	return p
}

func (p *ProgressModel) Wait() *ProgressModel {
	p.wait.Wait()
	return p
}

func (p *ProgressModel) Waiting() bool {
	return p.waitCount > 0
}

func (p *ProgressModel) ShowError(err error, msg ...string) *ProgressModel {
	p.Clear().
		SetWait().
		Show(msg...).
		Show("The following error has occurred.").
		ShowDone(err.Error()).
		Wait()
	return p
}

func (p *ProgressModel) Show(newLine ...string) *ProgressModel {
	return p.doShow(false, newLine...)
}

func (p *ProgressModel) ShowDone(newLine ...string) *ProgressModel {
	return p.doShow(true, newLine...)
}

func (p *ProgressModel) doShow(done bool, newLine ...string) *ProgressModel {
	if p.modal != nil {
		p.modal.Hide()
	}

	p.progressList = append(p.progressList, newLine...)
	p.modal = widget.NewModalPopUp(
		widget.NewRichTextWithText("Nothing to see here, yet..."),
		p.window.Canvas(),
	)
	p.modal.Content = p.progressBoxInfi(done)
	p.modal.Show()

	return p
}

func (p *ProgressModel) Clear() *ProgressModel {
	p.progressList = make([]string, 0)
	return p
}

func (p *ProgressModel) ShowProgress(minVal, maxVal float64,
	textFormatter func() string, newLine ...string) *ProgressModel {
	if p.modal != nil {
		p.modal.Hide()
	}
	p.progressList = append(p.progressList, newLine...)
	p.modal = widget.NewModalPopUp(
		widget.NewRichTextWithText("Nothing to see here, yet..."),
		p.window.Canvas(),
	)
	p.modal.Content = p.progressBox(minVal, maxVal, textFormatter)
	p.modal.Show()
	return p
}

func (p *ProgressModel) UpdateProgress(values ...interface{}) *ProgressModel {
	for _, value := range values {
		switch v := value.(type) {
		case float64:
			p.progressBar.SetValue(v)
		case string:
			last := len(p.progressList) - 1
			p.progressList[last] = v
			p.modal.Content.(*fyne.Container).Objects[last] = p.makeProgressList()[last]
			p.modal.Refresh()
		default:
		}
	}
	return p
}

func (p *ProgressModel) Hide() *ProgressModel {
	if p.modal != nil {
		p.modal.Hide()
	}
	return p
}

func (p *ProgressModel) Bar(show bool) *ProgressModel {
	p.infiniteProgressBar = show
	return p
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

func (p *ProgressModel) progressBox(minVal, maxVal float64, textFormatter func() string) *fyne.Container {
	p.progressBar = &widget.ProgressBar{
		Min:   minVal,
		Max:   maxVal,
		Value: 0,
	}
	if textFormatter != nil {
		p.progressBar.TextFormatter = textFormatter
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
