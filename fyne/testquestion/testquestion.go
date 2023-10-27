package testquestion

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/abaskin/testparts"
	"github.com/daichi-m/go18ds/maps/linkedhashmap"
)

type TestQuestion struct {
	question      *testparts.GormQuestion
	allocTime     time.Duration
	options       *LabelRadioGroup
	next          *widget.Button
	answer        string
	correctAnswer string
}

func NewTestQuestion(question *testparts.GormQuestion,
	allocTime time.Duration, next *widget.Button) *TestQuestion {
	q := &TestQuestion{
		question:  question,
		allocTime: allocTime,
		next:      next,
	}
	q.options = NewLabelRadioGroup(q, []string{}, nil)
	for _, choice := range q.question.Choices {
		q.options.Append(choice.Choice)
		if choice.Answer {
			q.correctAnswer = choice.Choice
		}
	}
	q.options.onChanged = func(s string) {
		q.answer = s
	}
	return q
}

func (q *TestQuestion) Ask(countDown *widget.ProgressBar, content *fyne.Container) {
	*content = *q.Show()

	qDone := make(chan bool)
	q.next.OnTapped = func() { qDone <- true }

	countDown.Min = 0
	countDown.Max = float64(q.allocTime.Milliseconds())
	countDown.TextFormatter = func() string { return "" }

	timeRemain := q.allocTime
	qTicker := testparts.NewTicker(time.Millisecond, nil,
		func() {
			countDown.Value = float64(timeRemain.Milliseconds())
			countDown.Refresh()
			if timeRemain > 0 {
				timeRemain -= time.Millisecond
			}
		})
	qTimer := time.AfterFunc(q.allocTime, func() { qDone <- true })

	<-qDone
	qTimer.Stop()
	qTicker.Stop()

	q.next.OnTapped = nil
}

func (q *TestQuestion) Show() *fyne.Container {
	return container.NewVBox(
		layout.NewSpacer(),
		&widget.Label{
			Text:      q.question.Question,
			Wrapping:  fyne.TextWrapWord,
			Alignment: fyne.TextAlignCenter,
			TextStyle: fyne.TextStyle{
				Bold: true,
			},
		},
		container.NewCenter(q.options),
		layout.NewSpacer(),
	)
}

type LabelRadioGroup struct {
	question  *TestQuestion
	options   *linkedhashmap.Map[string, string]
	onChanged func(string)
	widget.RadioGroup
}

func NewLabelRadioGroup(question *TestQuestion, options []string,
	change func(string)) *LabelRadioGroup {
	newRG := &LabelRadioGroup{
		question:  question,
		options:   linkedhashmap.New[string, string](),
		onChanged: change,
	}
	newRG.ExtendBaseWidget(newRG)
	newRG.OnChanged = func(s string) {
		if newRG.onChanged != nil {
			newRG.onChanged(newRG.Answer())
		}
	}
	newRG.Append(options...)

	return newRG
}

func (l *LabelRadioGroup) Append(options ...string) {
	for _, opt := range options {
		label := fmt.Sprintf("%c.  %s", 'A'+l.options.Size(), opt)
		l.options.Put(label, opt)
		l.Append(label)
	}
}

func (l *LabelRadioGroup) Answer() string {
	answer, _ := l.options.Get(l.Selected)
	return answer
}

func (l *LabelRadioGroup) SetSelect(keyPressed fyne.KeyName) {
	if keyPressed == fyne.KeyReturn {
		test.Tap(l.question.next)
	}
	if opt, _ := l.options.Find(func(key, value string) bool {
		return strings.HasPrefix(key, string(keyPressed))
	}); opt != "" {
		l.SetSelected(opt)
	}
}
