package testquestion

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/abaskin/testparts"
	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/daichi-m/go18ds/maps/linkedhashmap"
	"github.com/hako/durafmt"
)

type TestQuestion struct {
	Question      *testparts.GormQuestion
	Options       *LabelRadioGroup
	OptionList    *arraylist.List[*ClickText]
	Done          bool
	Resolution    time.Duration
	CorrectAnswer *string
	Answer        string
	Feedback      *string
	allocTime     time.Duration
	next          *widget.Button
}

func NewTestQuestion(question *testparts.GormQuestion, allocTime time.Duration,
	next *widget.Button) *TestQuestion {
	q := &TestQuestion{
		Question:   question,
		OptionList: arraylist.New[*ClickText](),
		Done:       false,
		Resolution: time.Millisecond,
		allocTime:  allocTime,
		next:       next,
	}
	q.Options = NewLabelRadioGroup(q, []string{}, nil)
	for i, choice := range q.Question.Choices {
		q.Options.Add(choice.Choice)
		q.OptionList.Add(NewClickText(q, choice))
		if choice.Answer {
			q.CorrectAnswer = &q.Question.Choices[i].Choice
		}
	}
	q.Options.onChanged = func(choice string) {
		q.Answer = choice
		q.SetFeedback()
	}
	return q
}

func (q *TestQuestion) Ask(countDown *widget.ProgressBar, content *fyne.Container,
	fullWindow *fyne.Container, qDone *sync.WaitGroup) {
	*content = *q.Show()
	content.Refresh()
	fullWindow.Refresh()

	if q.next != nil {
		q.next.OnTapped = func() { qDone.Done() }
	}

	timeRemain := q.allocTime
	countDown.Min = 0
	countDown.Max = float64(q.allocTime.Milliseconds())
	countDown.TextFormatter =
		func() string {
			return durafmt.Parse(timeRemain.Truncate(time.Second)).String()
		}

	qTicker := testparts.NewTicker(q.Resolution, nil,
		func() {
			countDown.SetValue(float64(timeRemain.Milliseconds()))
			if timeRemain > 0 {
				timeRemain -= q.Resolution
			}
		})

	qDone.Wait()
	qTicker.Stop()

	if q.next != nil {
		q.next.OnTapped = nil
	}
}

func (q *TestQuestion) Show() *fyne.Container {
	showContainer := container.NewVBox(
		&widget.RichText{
			Wrapping:   fyne.TextWrapWord,
			Scroll:     container.ScrollNone,
			Truncation: fyne.TextTruncateOff,
			Segments: []widget.RichTextSegment{
				&widget.TextSegment{
					Text: q.Question.Question,
					Style: widget.RichTextStyle{
						TextStyle: fyne.TextStyle{
							Bold: true,
						},
						Alignment: fyne.TextAlignCenter,
						SizeName:  "QuestionFontSize",
						ColorName: "QuestionColor",
					},
				},
			},
		},
	)

	q.OptionList.Each(func(_ int, opt *ClickText) {
		showContainer.Add(opt)
	})

	return showContainer
}

func (q *TestQuestion) Correct() bool {
	return q.Answer == *q.CorrectAnswer
}

func (q *TestQuestion) AnswerID() uint {
	for _, c := range q.Question.Choices {
		if q.Answer == c.Choice {
			return c.ID
		}
	}
	return 0
}

func (q *TestQuestion) SetFeedback() {
	q.Feedback = &q.Question.Choices[slices.IndexFunc(q.Question.Choices,
		func(c testparts.GormQuestionChoice) bool {
			return q.Answer == c.Choice
		})].Feedback
}

func (q *TestQuestion) SetSelect(keyPressed fyne.KeyName) {
	if keyPressed == fyne.KeyReturn {
		test.Tap(q.next)
	}
	q.OptionList.Find(
		func(_ int, c *ClickText) bool {
			if strings.HasPrefix(
				c.RichText.Segments[0].(*widget.TextSegment).Text, string(keyPressed)) {
				test.Tap(c)
			}
			return false
		})
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
	newRG.Add(options...)

	return newRG
}

func (l *LabelRadioGroup) Add(options ...string) {
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

func formatChoice(choice string) *widget.RichText {
	return &widget.RichText{
		Wrapping:   fyne.TextWrapWord,
		Scroll:     container.ScrollNone,
		Truncation: fyne.TextTruncateOff,
		Segments: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: choice,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{
						Bold: true,
					},
					Alignment: fyne.TextAlignCenter,
					SizeName:  "OptionFontSize",
					ColorName: "OptionColor",
				},
			},
		},
	}
}

type ClickText struct {
	question *TestQuestion
	choice   testparts.GormQuestionChoice
	*widget.RichText
}

func NewClickText(question *TestQuestion, choice testparts.GormQuestionChoice) *ClickText {
	ct := &ClickText{
		RichText: formatChoice(
			fmt.Sprintf("%c.  %s", 'A'+question.OptionList.Size(), choice.Choice)),
		question: question,
		choice:   choice,
	}
	ct.ExtendBaseWidget(ct)
	return ct
}

func (ct *ClickText) Tapped(p *fyne.PointEvent) {
	ct.question.Answer = ct.choice.Choice
	ct.question.SetFeedback()
	ct.SetSelected()
}

func (ct *ClickText) SetSelected() {
	ct.question.OptionList.Each(
		func(_ int, c *ClickText) {
			style := &c.RichText.Segments[0].(*widget.TextSegment).Style
			style.ColorName = "OptionColor"
			style.TextStyle.Italic = false
			if c == ct {
				style.ColorName = "OptionColorSelected"
				style.TextStyle.Italic = true
			}
			c.RichText.Refresh()
		})
}

func (ct *ClickText) Correct() bool {
	return ct.choice.Answer
}
