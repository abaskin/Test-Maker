package quizgui

import (
	"image/color"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/abaskin/testparts"
	icon "github.com/abaskin/testparts/fyne/icons"
	"github.com/abaskin/testparts/fyne/progressmodel"
	"github.com/abaskin/testparts/quizconfig"
)

func SetupGUI(gui *quizconfig.Gui_t) error {
	gui.MyApp.Settings().SetTheme(&quizconfig.QuizTheme{})
	gui.MyWindow.Resize(fyne.NewSize(1024, 768))
	gui.MyWindow.SetMaster()
	// gui.MyWindow.SetFullScreen(true)
	gui.MyWindow.SetCloseIntercept(func() {
		log.Println("SetCloseIntercept, Main window closing")
		gui.MyWindow.Close()
	})

	gui.StatusModel = progressmodel.NewProgressModel(gui.MyWindow)

	gui.Text = make(map[string]*canvas.Text)

	gui.AddText("clock", 25, &color.RGBA{0, 0, 255, 255}, fyne.TextAlignCenter)
	gui.ClockTicker = testparts.NewTicker(time.Second,
		func() {
			// not used
		},
		func() {
			gui.Text["clock"].Text = func() string {
				currentTime := time.Now()
				return currentTime.Format("3:04:05 pm")
			}()
			gui.Text["clock"].Refresh()
		})

	gui.ExitButton = widget.NewButton("Exit",
		func() {
			// check to see if we can close
			gui.MyWindow.Close()
		})
	gui.AnswerButton = widget.NewButton("Answer", nil)

	gui.Scores = &fyne.Container{}

	gui.CenterContent = &fyne.Container{}
	gui.QuestionCountdown = &widget.ProgressBar{}

	gui.Scores = &fyne.Container{}

	gui.Avatar = &canvas.Image{
		ScaleMode: canvas.ImageScaleFastest,
		FillMode:  canvas.ImageFillContain,
	}
	gui.Avatar.SetMinSize(fyne.NewSize(60, 60))

	gui.HeaderLeft = &fyne.Container{}

	gui.HeaderRight = &canvas.Text{
		Text:     "",
		TextSize: 30,
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
		Alignment: fyne.TextAlignTrailing,
	}

	gui.FullWindow = container.NewBorder(
		container.NewVBox(
			container.NewGridWithColumns(3,
				gui.HeaderLeft,
				gui.QuestionCountdown,
				gui.HeaderRight,
			),
			gui.Scores,
		),
		container.NewGridWithColumns(3,
			gui.ExitButton,
			gui.Text["clock"],
			gui.AnswerButton,
		),
		nil, nil,
		gui.CenterContent,
	)

	gui.MyWindow.SetContent(gui.FullWindow)

	return nil
}

func scoreTextFmt(text string, textColor color.RGBA) *canvas.Text {
	return &canvas.Text{
		Text:      text,
		TextSize:  30,
		TextStyle: fyne.TextStyle{Bold: true},
		Alignment: fyne.TextAlignCenter,
		Color:     FromRGBA(textColor),
	}
}

func UpdateScores(clientList []*quizconfig.Client_t, scores *fyne.Container) {
	scores.RemoveAll()

	for _, client := range clientList {
		avatar := &canvas.Image{
			ScaleMode: canvas.ImageScaleFastest,
			FillMode:  canvas.ImageFillContain,
			Resource:  icon.Avatars[client.AvatarName],
		}
		avatar.SetMinSize(fyne.NewSize(60, 60))

		scores.Add(
			container.NewHBox(
				avatar,
				container.NewVBox(
					scoreTextFmt(client.Name, client.Color),
					scoreTextFmt(strconv.Itoa(int(client.Score)), client.Color),
				),
			),
		)
	}

	scores.Refresh()
}

func ToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

func FromRGBA(c color.RGBA) color.Color {
	return color.RGBA{c.R, c.G, c.B, c.A}
}
