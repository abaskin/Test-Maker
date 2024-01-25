package datafield

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/golang-module/carbon/v2"
	"github.com/monitor1379/yagods/sets/treeset"
	"github.com/samber/lo"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// ---- DataButton ----

type DataButton[T any] struct {
	*widget.Button
	Data     T
	OnTapped func(*DataButton[T]) `json:"-"`
}

func (btn *DataButton[T]) Tapped(*fyne.PointEvent) {
	if btn.OnTapped != nil {
		btn.OnTapped(btn)
	}
}

func NewButton[T any](text string, icon fyne.Resource, data T,
	tapped func(*DataButton[T])) *DataButton[T] {
	btn := &DataButton[T]{
		OnTapped: tapped,
		Data:     data,
		Button: &widget.Button{
			Text: text,
			Icon: icon,
		},
	}

	btn.ExtendBaseWidget(btn)
	return btn
}

// ---- DataCheck ----

type DataCheck[T any] struct {
	*widget.Check
	Data     T
	OnTapped func(*DataCheck[T]) `json:"-"`
}

func (btn *DataCheck[T]) Tapped(*fyne.PointEvent) {
	if btn.OnTapped != nil {
		btn.OnTapped(btn)
	}
}

func NewCheck[T any](text string, state bool, data T,
	tapped func(*DataCheck[T])) *DataCheck[T] {
	btn := &DataCheck[T]{
		OnTapped: tapped,
		Data:     data,
	}
	btn.Check = &widget.Check{
		Text: text,
	}
	btn.Check.SetChecked(state)

	btn.ExtendBaseWidget(btn)
	return btn
}

// ---- DataSelect ----

type SelectData[K comparable, V any] struct {
	*orderedmap.OrderedMap[K, V]
}

func (sd *SelectData[K, V]) Size() int {
	return sd.Len()
}

func (sd *SelectData[K, V]) Keys() []K {
	keys, _ := sd.Parts()
	return keys
}

func (sd *SelectData[K, V]) Values() []V {
	_, values := sd.Parts()
	return values
}

func (sd *SelectData[K, V]) Parts() ([]K, []V) {
	keys := make([]K, 0)
	values := make([]V, 0)
	for pair := sd.Oldest(); pair != nil; pair = pair.Next() {
		keys = append(keys, pair.Key)
		values = append(values, pair.Value)
	}
	return keys, values
}

func (sd *SelectData[K, V]) Each(eachFunc func(K, V)) {
	for pair := sd.Oldest(); pair != nil; pair = pair.Next() {
		eachFunc(pair.Key, pair.Value)
	}
}

func NewSelectData[K comparable, V any]() *SelectData[K, V] {
	return &SelectData[K, V]{
		OrderedMap: orderedmap.New[K, V](),
	}
}

type DataSelect[T any] struct {
	*widget.Select
	Data *SelectData[string, T]
}

func NewSelect[T any](data *SelectData[string, T]) *DataSelect[T] {
	sel := &DataSelect[T]{
		Data: data,
		Select: &widget.Select{
			Options: data.Keys(),
		},
	}

	sel.ExtendBaseWidget(sel)
	return sel
}

func (sel *DataSelect[T]) Value() T {
	value, found := sel.Data.Get(sel.Selected)
	return lo.Ternary(found, value, *new(T))
}

// ---- DataEntry ----

type DataEntry[T any] struct {
	*widget.Entry
	Data    T
	Changed func(string, *DataEntry[T]) `json:"-"`
}

func NewEntry[T any](text string, data T,
	change func(string, *DataEntry[T])) *DataEntry[T] {
	fld := &DataEntry[T]{
		Changed: change,
		Data:    data,
		Entry: &widget.Entry{
			Text:        text,
			PlaceHolder: "      ",
		},
	}
	fld.OnChanged = func(s string) {
		if fld.Changed != nil {
			fld.Changed(s, fld)
		}
	}

	fld.ExtendBaseWidget(fld)
	return fld
}

// func (fld *DataEntry[T]) OnChanged(s string) {
// 	if fld.Changed != nil {
// 		fld.Changed(s, fld)
// 	}
// }

// ---- TimePicker ----

type TimePicker struct {
	time   *carbon.Carbon
	hour   *widget.SelectEntry
	minute *widget.SelectEntry
	ampm   *widget.SelectEntry
}

func NewTimePicker(t *carbon.Carbon, tl *TimeLabel) *TimePicker {
	hourStr := make([]string, 0)
	for n := 1; n < 13; n++ {
		hourStr = append(hourStr, strconv.Itoa(n))
	}
	tp := &TimePicker{
		time:   t,
		hour:   widget.NewSelectEntry(hourStr),
		minute: widget.NewSelectEntry([]string{}),
		ampm:   widget.NewSelectEntry([]string{"AM", "PM"}),
	}

	tp.hour.OnChanged = func(hourStr string) {
		if tp.time.Format("g") != hourStr {
			hour, _ := strconv.Atoi(hourStr)
			*tp.time = tp.time.SetHour(hour)
			tl.Update()
		}
	}

	tp.minute.OnChanged = func(minuteStr string) {
		if tp.time.Format("i") != minuteStr {
			minute, _ := strconv.Atoi(minuteStr)
			*tp.time = tp.time.SetMinute(minute)
			tl.Update()
		}
	}

	tp.ampm.OnChanged = func(ap string) {
		if tp.time.Format("A") != ap {
			switch ap {
			case "AM":
				*tp.time = tp.time.SubHours(12)
			case "PM":
				*tp.time = tp.time.AddHours(12)
			default:
			}
			tl.Update()
		}
	}

	minutes := treeset.NewWithStringComparator("00", "15", "30", "45")
	if !minutes.Contains(tp.time.Format("i")) {
		minutes.Add(tp.time.Format("i"))
	}
	tp.minute.SetOptions(minutes.Values())

	tp.hour.SetText(tp.time.Format("g"))
	tp.minute.SetText(tp.time.Format("i"))
	tp.ampm.SetText(tp.time.Format("A"))

	return tp
}

func (tp *TimePicker) Render() *fyne.Container {
	if tp.time.IsZero() {
		return container.NewWithoutLayout()
	}
	return container.NewHBox(
		tp.hour,
		widget.NewLabel(":"),
		tp.minute,
		layout.NewSpacer(),
		tp.ampm,
	)
}

// ---- TimeLabel ----

type TimeLabel struct {
	*widget.Label
	time *carbon.Carbon
}

func NewTimeLabel(t *carbon.Carbon) *TimeLabel {
	tl := &TimeLabel{
		time: t,
		Label: widget.NewLabelWithStyle(
			"",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
	}
	tl.Update()

	tl.ExtendBaseWidget(tl)
	return tl
}

func (tl *TimeLabel) Update() {
	if !tl.time.IsZero() {
		tl.Label.SetText(tl.time.ToDayDateTimeString())
		tl.Label.Refresh()
	}
}
