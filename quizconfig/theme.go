package quizconfig

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	themeX "fyne.io/x/fyne/theme"
)

type QuizTheme struct{}

var _ fyne.Theme = (*QuizTheme)(nil)

var themeColorTable map[fyne.ThemeColorName]color.Color
var themeSizeTable map[fyne.ThemeSizeName]float32

func init() {
	themeColorTable = map[fyne.ThemeColorName]color.Color{
		"QuestionColor":       color.NRGBA{R: 0, G: 0, B: 255, A: 255},
		"OptionColor":         color.NRGBA{R: 0, G: 255, B: 0, A: 255},
		"OptionColorSelected": color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}

	themeSizeTable = map[fyne.ThemeSizeName]float32{
		"QuestionFontSize": 40,
		"OptionFontSize":   30,
	}
}

func (m QuizTheme) AddColor(name fyne.ThemeColorName, color color.Color) {
	themeColorTable[name] = color
}

func (m QuizTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case theme.VariantLight:
		if color, found := themeColorTable[name]; found {
			return color
		}
		switch name {
		case "OptionBGColor":
			return themeX.AdwaitaTheme().Color(theme.ColorNamePrimary, variant)
		default:
			return themeX.AdwaitaTheme().Color(name, variant)
		}
	case theme.VariantDark:
		switch name {
		default:
			return themeX.AdwaitaTheme().Color(name, variant)
		}
	default:
		return themeX.AdwaitaTheme().Color(name, variant)
	}
}

func (m QuizTheme) Font(style fyne.TextStyle) fyne.Resource {
	switch {
	// case style.Symbol:
	// 	return icon.NotoCOLRv1EmojicompatTtf
	default:
		return theme.DefaultTheme().Font(style)
	}
}

func (m QuizTheme) AddSize(name fyne.ThemeSizeName, size float32) {
	themeSizeTable[name] = size
}

func (m QuizTheme) Size(name fyne.ThemeSizeName) float32 {
	if size, found := themeSizeTable[name]; found {
		return size
	}
	return themeX.AdwaitaTheme().Size(name)
}

func (m QuizTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return themeX.AdwaitaTheme().Icon(name)
}
