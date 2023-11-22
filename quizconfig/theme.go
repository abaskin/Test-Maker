package quizconfig

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	themeX "fyne.io/x/fyne/theme"
)

type QuizTheme struct{}

var _ fyne.Theme = (*QuizTheme)(nil)

func (m QuizTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case theme.VariantLight:
		switch name {
		case "QuestionColor":
			return color.NRGBA{R: 0, G: 0, B: 255, A: 255}
		case "OptionColor":
			return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
		case "OptionColorSelected":
			return color.NRGBA{R: 255, G: 0, B: 0, A: 255}
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

func (m QuizTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case "QuestionFontSize":
		return 40
	case "OptionFontSize":
		return 30
	default:
		return themeX.AdwaitaTheme().Size(name)
	}
}

func (m QuizTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return themeX.AdwaitaTheme().Icon(name)
}
