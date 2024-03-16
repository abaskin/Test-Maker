package quizconfig

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/abaskin/testparts"
	"github.com/abaskin/testparts/aiclient"
	"github.com/abaskin/testparts/fyne/progressmodel"
	"github.com/daichi-m/go18ds/maps/linkedhashmap"
	"github.com/ebitengine/oto/v3"
	"github.com/looplab/fsm"
)

const (
	QoS = 1

	BrokerUrl     = "mqtt://broker.hivemq.com:1883"
	topicBase     = "/1Z0VReh9ReJ4Vx4NHAPoX16sIFPzlcrt"
	QuestionTopic = topicBase + "/sWnLDrZuuEQMrRYRrEJH09p9JVWoI4LM"
	ClientTopic   = topicBase + "/hTrun1EOZXxcVZpfPsvVU0FI11iqCAZE"

	MaxDatagramSize  = 8192
	MultiCastAddress = "239.192.0.10:1555"

	WaitTime = 3 * time.Second
)

type ClientData_t struct {
	Action                              Action_t
	Question                            *testparts.GormQuestion
	QuestionTime                        time.Duration
	QuestionNum, Name, UUID, AvatarName string
	Scores                              []*Client_t
	Color                               color.RGBA
	QuestionID                          uint
	Correct                             bool
}

type Gui_t struct {
	MyApp                     fyne.App
	MyWindow                  fyne.Window
	HeaderLeft                *fyne.Container
	QuestionCountdown         *widget.ProgressBar
	Scores                    *fyne.Container
	FullWindow, CenterContent *fyne.Container
	Text                      map[string]*canvas.Text
	HeaderRight               *canvas.Text
	StatusModel               *progressmodel.ProgressModel
	AnswerButton              *widget.Button
	ExitButton                *widget.Button
	Avatar                    *canvas.Image
	ClockTicker               *testparts.Ticker
}

type State_t struct {
	Transport              Transport_t
	Class                  *testparts.GormClass
	ClientWait             *sync.WaitGroup
	ClientList             *linkedhashmap.Map[string, *Client_t]
	RecvChan, ActionChan   chan *ClientData_t
	QuestionChan           chan *testparts.GormQuestion
	Name, UUID, AvatarName string
	Color                  color.Color
	Result                 map[Answer_t]*Result_t
	Player                 *oto.Player
	AiModel                aiclient.AiModel
}

type Client_t struct {
	Name, UUID, AvatarName string
	Score                  uint
	Color                  color.RGBA
}

type Result_t struct {
	Text string
	Icon *canvas.Image
}

type Transport_t interface {
	Send(data *ClientData_t, topic string) error
	SendNoop() error
	ShutDown()
}

type Action_t uint8

const (
	ActionUndefined Action_t = iota
	ActionHello
	ActionReHello
	ActionStart
	ActionQuestion
	ActionAnswer
	ActionAnswered
	ActionDone
	ActionNoop
)

func (a Action_t) String() string {
	actionString := map[Action_t]string{
		ActionUndefined: "Undefined",
		ActionHello:     "Hello",
		ActionReHello:   "ReHello",
		ActionStart:     "Start",
		ActionQuestion:  "Question",
		ActionAnswer:    "Answer",
		ActionAnswered:  "Answered",
		ActionDone:      "Done",
		ActionNoop:      "Noop",
	}
	return actionString[a]
}

type Answer_t uint8

const (
	AnswerUndefined Answer_t = iota
	AnswerCorrect
	AnswerWrong
	AnswerTooLate
)

func GetMetaObject[T any](fsm *fsm.FSM, key string) *T {
	obj, found := fsm.Metadata(key)
	if found {
		return obj.(*T)
	}
	return nil
}

func (gui *Gui_t) AddText(name string, size float32, clr color.Color,
	align fyne.TextAlign) {
	gui.Text[name] =
		&canvas.Text{
			Text:     "",
			TextSize: size,
			TextStyle: fyne.TextStyle{
				Bold: true,
			},
			Color:     clr,
			Alignment: align,
		}
}
