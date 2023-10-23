package fancybutton

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FancyButton struct {
	*widget.Card
	Data        interface{}
	Canvas      *canvas.Rectangle                    `json:"-"`
	OnTapped    func(*fyne.PointEvent, *FancyButton) `json:"-"`
	OnTappedSec func(*fyne.PointEvent, *FancyButton) `json:"-"`
	OnDragged   func(*fyne.DragEvent, *FancyButton)  `json:"-"`
	OnDragEnd   func(*FancyButton)                   `json:"-"`
	DragEvent   fyne.DragEvent                       `json:"-"`
}

func (fb *FancyButton) Dragged(de *fyne.DragEvent) {
	fb.DragEvent = *de
	if fb.OnDragged != nil {
		fb.OnDragged(de, fb)
	}
}

func (fb *FancyButton) DragEnd() {
	if fb.OnDragged != nil {
		fb.OnDragEnd(fb)
	}
}

func (fb *FancyButton) Tapped(p *fyne.PointEvent) {
	if fb.OnTapped != nil {
		fb.OnTapped(p, fb)
	}
}

func (fb *FancyButton) TappedSecondary(p *fyne.PointEvent) {
	if fb.OnTappedSec != nil {
		fb.OnTappedSec(p, fb)
	}
}

func New(content *fyne.Container, btnColor color.NRGBA,
	tapped func(*fyne.PointEvent, *FancyButton), data interface{}) *FancyButton {
	fb := &FancyButton{
		OnTapped: tapped,
		Data:     data,
	}
	fb.Canvas = canvas.NewRectangle(btnColor)
	fb.Card = widget.NewCard(
		"", "",
		container.NewMax(
			fb.Canvas,
			content,
		),
	)
	fb.ExtendBaseWidget(fb)
	return fb
}
