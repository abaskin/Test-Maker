package mapselect

import (
	"fyne.io/fyne/v2/widget"
	"github.com/daichi-m/go18ds/maps/linkedhashmap"
)

type MapSelect struct {
	options *linkedhashmap.Map[string, interface{}]
	widget.Select
}

func NewMapSelect(onChanged func(string)) *MapSelect {
	sel := &MapSelect{
		options: linkedhashmap.New[string, interface{}](),
	}
	sel.ExtendBaseWidget(sel)
	sel.OnChanged = onChanged
	return sel
}

func (m *MapSelect) Add(key string, value interface{}) {
	m.options.Put(key, value)
	m.Options = append(m.Options, key)
}

func (m *MapSelect) Key() interface{} {
	key, _ := m.options.Get(m.Selected)
	return key
}

func (m *MapSelect) Clear() {
	m.options.Clear()
	m.Options = []string{}
}
