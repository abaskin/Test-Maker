package testparts

import (
	"encoding/json"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/daichi-m/go18ds/maps/treemap"
)

const ZERO_WIDTH_SPACE = "\u200b"

type PathStrSt struct {
	TestFile, Outdir, Workdir, Assetdir, Templatepath, ImportFile string
}

type FlagsSt struct {
	ShowAll, SaveTex, CreateDistro, CreateRtf, CreateForm, CreatePDF, DBImport,
	ContinuousNumbering, ImportClass, ImportSession bool
}

type TestSt struct {
	TestJSON   TestJSONSt
	Sections   []SectionSt
	TestBundle *arraylist.List[*TestBundleSt]
	Template   []byte
}

type QuestgenQuestionSt struct {
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Choices  WordsSt `json:"distractors"`
}

type QuestionsSt struct {
	Answer   uint           `json:"answer"`
	NumCol   uint           `json:"numCol"`
	Used     uint           `json:"used"`
	Required bool           `json:"required"`
	Choices  WordsSt        `json:"choices"`
	Answers  WordsSt        `json:"answers"`
	Question NLStringSt     `json:"question"`
	Parts    NLStringListSt `json:"parts"`
}

type JSONSectionSt struct {
	Type             string        `json:"type"`
	SectionTitle     string        `json:"sectionTitle"`
	NumLines         string        `json:"numLines"`
	Title            string        `json:"title"`
	Points           uint          `json:"points"`
	NumQuest         uint          `json:"questionsOntest"`
	NumCol           uint          `json:"numCol"`
	AnswerLines      bool          `json:"answerLines"`
	QuizBox          bool          `json:"quizBox"`
	KeepOrder        bool          `json:"keepOrder"`
	AnswerText       WordsSt       `json:"answerText"`
	WordList         WordsSt       `json:"word-list"`
	Include          WordsSt       `json:"include"`
	IncludeQuestgen  WordsSt       `json:"includeQuestgen"`
	IncludeAiken     WordsSt       `json:"includeAiken"`
	Answers          WordsSt       `json:"answers"`
	ColumnHead       WordsSt       `json:"columnHead"`
	Instructions     NLStringSt    `json:"instructions"`
	FormInstructions NLStringSt    `json:"formInstructions"`
	Text             NLStringSt    `json:"text"`
	Words            WordDefMapSt  `json:"words"`
	Questions        QuestionSetSt `json:"questions"`
	AssetDir         string        `json:"-"`
}

type TestJSONSt struct {
	Subject      string          `json:"subject"`
	Grade        string          `json:"grade"`
	Title        string          `json:"title"`
	RTFTitle     string          `json:"rtfTitle"`
	School       string          `json:"school"`
	Logo         string          `json:"logo"`
	Date         string          `json:"date"`
	Time         uint            `json:"time"`
	NoKey        bool            `json:"noKey"`
	MinQuestions uint            `json:"minQuestions"`
	Students     WordsSt         `json:"students"`
	Classes      ClassMapSt      `json:"classes"`
	Sections     []JSONSectionSt `json:"sections"`
}

type ClassJSONSt struct {
	Name     string  `json:"name"`
	Students WordsSt `json:"students"`
}

type SectionHeadSt struct {
	Type             string
	SectionTitle     string
	Title            string
	NumLines         string
	Points           uint
	NumCol           uint
	NumQuest         uint
	AnswerLines      bool
	QuizBox          bool
	Quiz             bool
	KeepOrder        bool
	Instructions     string
	FormInstructions string
	Text             string
	ColumnHead       WordsSt
}

type SectionSt interface {
	Init(JSONSectionSt, uint)
	TestLatex(uint, *QuestNumSt) []string
	TestRTF(*RTFDoc, uint, *QuestNumSt)
	TestForm(*GoogleFormSt, uint) error
	AnswerLatex(bool, bool, uint, *QuestNumSt) []string
	DistribLatex(string, uint) []string
	GetHead() *SectionHeadSt
}

type WordMatchSt struct {
	SectionHeadSt
	*arraylist.List[WordListSt]
	WordDist     WordDistMapSt
	AllQuestions QuestionSetSt
	AllWords     *WordDefMapSt
}

type MultipleChoiceSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
}

type WordProblemSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
}

type CompQuestionsSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
}

type ReadingCompSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
}

type QuizSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
}

type CustomSt struct {
	SectionHeadSt
	Questions    QuestionListSt
	AllQuestions QuestionSetSt
	AnswerText   *arraylist.List[string]
	Answers      *arraylist.List[string]
}

type PassageCompletionSt struct {
	SectionHeadSt
	WordList  StringListSt
	Answers   StringListSt
	Questions QuestionListSt // not used
}

type TestHeadSt struct {
	Subject  string
	Grade    string
	School   string
	Title    string
	RTFTitle string
	Logo     string
	Date     string
	Points   uint
	Time     uint
	Quiz     bool
	NoKey    bool
	Dsn      string
}

type TestBundleSt struct {
	Student    string
	StudentNum uint
	*TestHeadSt
}

type TakerClass struct {
	Subject  string `json:"subject"`
	Sections []struct {
		Section  string   `json:"section"`
		Students []string `json:"students"`
	} `json:"sections"`
}

type TakerSession struct {
	Subject      string    `json:"subject"`
	Title        string    `json:"title"`
	Section      string    `json:"section"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	QuestionTime string    `json:"questionTime"`
}

type ClassMapSt struct {
	*treemap.Map[string, WordsSt]
}

func (cm *ClassMapSt) UnmarshalJSON(data []byte) error {
	classes := make([]*ClassJSONSt, 0)
	if err := json.Unmarshal(data, &classes); err != nil {
		return err
	}
	cm.Map = &treemap.Map[string, WordsSt]{}
	for _, class := range classes {
		cm.Map.Put(class.Name, class.Students)
	}
	return nil
}

type QuestionSetSt struct {
	*arraylist.List[*QuestionsSt]
}

func (qs *QuestionSetSt) UnmarshalJSON(data []byte) error {
	qSet := make([]*QuestionsSt, 0)
	if err := json.Unmarshal(data, &qSet); err != nil {
		return err
	}
	qs.List = arraylist.New(qSet...)
	return nil
}

func (qs *QuestionSetSt) fixMissing() {
	if qs.List == nil {
		qs.List = &arraylist.List[*QuestionsSt]{}
	}
}

type QuestionListSt struct {
	*arraylist.List[QuestionSetSt]
}

func (ql QuestionListSt) get(index uint) QuestionSetSt {
	if questions, found := ql.Get(int(index)); found {
		return questions
	}
	return QuestionSetSt{}
}

type WordDefSt struct {
	Word   NLStringSt
	Def    NLStringSt
	Answer NLStringSt
}

type WordListSt struct {
	*arraylist.List[*WordDefSt]
}

func (wm WordMatchSt) getWords(student uint) NLStringListSt {
	words := NLStringListSt{List: arraylist.New[NLStringSt]()}
	if wordList, found := wm.Get(int(student)); found {
		wordList.Each(
			func(_ int, w *WordDefSt) {
				words.Add(w.Word)
			},
		)
	}
	return words
}

type WordDefMapSt struct {
	*treemap.Map[NLStringSt, NLStringSt]
}

func newWordDefMap() *treemap.Map[NLStringSt, NLStringSt] {
	return treemap.NewWith[NLStringSt, NLStringSt](
		func(a, b NLStringSt) int {
			switch {
			case a.string > b.string:
				return 1
			case a.string < b.string:
				return -1
			default:
				return 0
			}
		},
	)
}

func (wdm *WordDefMapSt) UnmarshalJSON(data []byte) error {
	wSet := make(map[string]string, 0)
	if err := json.Unmarshal(data, &wSet); err != nil {
		return err
	}
	wdm.Map = newWordDefMap()
	re, _ := regexp.Compile(`^See[\s\w]+\.`)
	for k, v := range wSet {
		if matched := re.MatchString(v); !matched {
			wdm.Put(
				NLStringSt{strings.ReplaceAll(k, ZERO_WIDTH_SPACE, ``)},
				NLStringSt{strings.ReplaceAll(v, ZERO_WIDTH_SPACE, ``)},
			)
		}
	}
	return nil
}

func (wdm *WordDefMapSt) get(word NLStringSt) NLStringSt {
	if def, found := wdm.Get(word); found {
		return def
	}
	return NLStringSt{""}
}

func (wdm *WordDefMapSt) fixMissing() {
	if wdm.Map == nil {
		wdm.Map = &treemap.Map[NLStringSt, NLStringSt]{}
	}
}

type WordDistMapSt struct {
	*treemap.Map[NLStringSt, uint]
}

func newWordDistMap() *treemap.Map[NLStringSt, uint] {
	return treemap.NewWith[NLStringSt, uint](
		func(a, b NLStringSt) int {
			switch {
			case a.string > b.string:
				return 1
			case a.string < b.string:
				return -1
			default:
				return 0
			}
		},
	)
}

type WordsSt struct {
	*arraylist.List[string]
}

func (w *WordsSt) UnmarshalJSON(data []byte) error {
	w.List = arraylist.New[string]()
	if err := w.FromJSON(data); err != nil {
		return err
	}
	return nil
}

func (w *WordsSt) fixMissing() WordsSt {
	if w.List == nil {
		w.List = arraylist.New[string]()
	}
	return *w
}

type NLStringSt struct {
	string
}

func (nls *NLStringSt) UnmarshalJSON(data []byte) error {
	wSet := make([]string, 0)
	if err := json.Unmarshal(data, &wSet); err != nil {
		return err
	}
	nls.string = strings.ReplaceAll(strings.Join(wSet, "\n"), "\u200b", ``)

	return nil
}

func (nls *NLStringSt) rtfString() string {
	outStr := strings.ReplaceAll(nls.string, "\n", " ")
	outStr = strings.ReplaceAll(outStr, "\u2019", `\'92`)
	outStr = strings.ReplaceAll(outStr, "\u201c", `\'93`)
	outStr = strings.ReplaceAll(outStr, "\u201d", `\'94`)
	outStr = strings.ReplaceAll(outStr, "\u2013", `\'96`)
	outStr = strings.ReplaceAll(outStr, `\fillin\`, "________")
	return outStr
}

func (nls *NLStringSt) CleanString() string {
	outStr := strings.ReplaceAll(nls.string, "\n", " ")
	outStr = strings.ReplaceAll(outStr, `\fillin\`, "________")
	return outStr
}

type NLStringListSt struct {
	*arraylist.List[NLStringSt]
}

func (nls *NLStringListSt) UnmarshalJSON(data []byte) error {
	wSet := make([][]string, 0)
	if err := json.Unmarshal(data, &wSet); err != nil {
		return err
	}
	nls.List = arraylist.New[NLStringSt]()
	for _, strs := range wSet {
		nls.Add(NLStringSt{strings.Join(strs, "\n")})
	}
	return nil
}

func (nls *NLStringListSt) getRandom() string {
	s, _ := nls.Get(rand.Intn(nls.Size()))
	return s.string
}

func (nls *NLStringListSt) fixMissing() NLStringListSt {
	if nls.List == nil {
		nls.List = arraylist.New[NLStringSt]()
	}
	return *nls
}

func (nls *NLStringListSt) join(sep string) string {
	outStr := make([]string, 0)
	nls.Each(
		func(_ int, value NLStringSt) {
			outStr = append(outStr, value.string)
		},
	)
	return strings.Join(outStr, sep)
}

func (nls *NLStringListSt) values() []string {
	outStr := make([]string, 0)
	nls.Each(
		func(_ int, value NLStringSt) {
			outStr = append(outStr, value.string)
		},
	)
	return outStr
}

func (wm WordMatchSt) getDefs(student uint) NLStringListSt {
	defs := NLStringListSt{List: arraylist.New[NLStringSt]()}
	if wordList, found := wm.Get(int(student)); found {
		wordList.Each(
			func(_ int, w *WordDefSt) {
				defs.Add(w.Def)
			},
		)
	}
	return defs
}

func (wm WordMatchSt) getAnswers(student uint) *arraylist.List[string] {
	answers := arraylist.New[string]()
	if wordList, found := wm.Get(int(student)); found {
		wordList.Each(
			func(_ int, w *WordDefSt) {
				answers.Add(w.Answer.string)
			},
		)
	}
	return answers
}

type StringListSt struct {
	*arraylist.List[NLStringListSt]
}

func (sl StringListSt) UnmarshalJSON(data []byte) error {
	wSet := make([][]string, 0)
	if err := json.Unmarshal(data, &wSet); err != nil {
		return err
	}
	sl.List = arraylist.New[NLStringListSt]()
	for _, words := range wSet {
		newList := NLStringListSt{List: arraylist.New[NLStringSt]()}
		for _, w := range words {
			newList.Add(NLStringSt{w})
		}
		sl.Add(newList)
	}
	return nil
}

func (sl StringListSt) get(index int) NLStringListSt {
	if questions, found := sl.Get(index); found {
		return questions
	}
	return NLStringListSt{List: arraylist.New[NLStringSt]()}
}

type QuestNumSt struct {
	number           uint32
	sectionNumbering bool
}

func MakeQuestNum(sectionNumbering bool) QuestNumSt {
	return QuestNumSt{
		number:           0,
		sectionNumbering: sectionNumbering,
	}
}

func (qn *QuestNumSt) NextNumber() uint32 {
	qn.number++
	return qn.number
}

func (qn *QuestNumSt) AddNumber(add uint32) {
	qn.number += add
}

func (qn *QuestNumSt) CurrentNumber() uint32 {
	return qn.number
}

func (qn *QuestNumSt) NewSection() {
	if qn.sectionNumbering {
		qn.number = 0
	}
}

func stringsUsing[T any](ss []T, transform func(T) string) []string {
	strOut := arraylist.New[string]()
	for _, value := range ss {
		strOut.Add(transform(value))
	}
	return strOut.Values()
}

// SequenceUsing generates slice in range using creator function
//
// There are 3 variations to generate:
//  1. [0, n).
//  2. [min, max).
//  3. [min, max) with step.
//
// if len(params) == 1 considered that will be returned slice between 0 and n,
// where n is the first param, [0, n).
// if len(params) == 2 considered that will be returned slice between min and max,
// where min is the first param, max is the second, [min, max).
// if len(params) == 3 considered that will be returned slice between min and max with step,
// where min is the first param, max is the second, step is the third one, [min, max) with step
func sequenceUsing[T any](ss []T, creator func(int) T, params ...int) []T {
	var seq = func(min, max, step int) (seq []T) {
		lenght := int(math.Round(float64(max-min) / float64(step)))
		if lenght < 1 {
			return
		}

		seq = make([]T, lenght)
		for i := 0; i < lenght; min += step {
			seq[i] = creator(min)
			i++
		}

		return seq
	}

	switch len(params) {
	case 1:
		return seq(0, params[0], 1)
	case 2:
		return seq(params[0], params[1], 1)
	case 3:
		return seq(params[0], params[1], params[2])
	default:
		return nil
	}
}

func ternary[T any](condition bool, ifOutput T, elseOutput T) T {
	switch condition {
	case true:
		return ifOutput
	default:
		return elseOutput
	}
}

func shuffleSlice[T any](items []T) []T {
	newItems := items
	rand.Shuffle(
		len(newItems),
		func(i, j int) {
			newItems[i], newItems[j] = newItems[j], newItems[i]
		},
	)
	return newItems
}
