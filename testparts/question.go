package testparts

import (
	"encoding/json"
	"fmt"
	"strings"

	aiken "github.com/aldinokemal/go-aiken"
	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/daichi-m/go18ds/sets/linkedhashset"
	"github.com/nwillc/genfuncs"
	"muzzammil.xyz/jsonc"
)

func (m *MultipleChoiceSt) Init(section JSONSectionSt, numTest uint) {
	m.SectionHeadSt = getSectionHead(section)
	m.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	m.AllQuestions = section.Questions
}

func (r *ReadingCompSt) Init(section JSONSectionSt, numTest uint) {
	r.SectionHeadSt = getSectionHead(section)
	r.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	r.AllQuestions = section.Questions
}

func (w *WordProblemSt) Init(section JSONSectionSt, numTest uint) {
	w.SectionHeadSt = getSectionHead(section)
	w.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	w.AllQuestions = section.Questions
}

func (q *QuizSt) Init(section JSONSectionSt, numTest uint) {
	q.SectionHeadSt = getSectionHead(section)
	q.Quiz = true
	q.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	q.AllQuestions = section.Questions
}

func (w *WordMatchSt) Init(section JSONSectionSt, numTest uint) {
	w.SectionHeadSt = getSectionHead(section)
	if w.ColumnHead.Empty() {
		w.ColumnHead.Add("Word", "Definition")
	}
	w.List = arraylist.New[WordListSt]()
	w.WordDist = WordDistMapSt{Map: newWordDistMap()}
	w.AllWords = &section.Words

	section.Words.Each(
		func(word, _ NLStringSt) {
			w.WordDist.Put(word, 0)
		},
	)

	for w.Size() < int(numTest) {
		wList := WordListSt{
			List: arraylist.New[*WordDefSt](),
		}

		w.WordDist.WordSet(section.NumQuest, numTest).Each(
			func(_ int, word NLStringSt) {
				wList.Add(&WordDefSt{Word: word})
			},
		)

		defines := make([]NLStringSt, 0)
		wList.Each(
			func(_ int, wd *WordDefSt) {
				defines = append(defines, section.Words.get(wd.Word))
			},
		)
		defines = shuffleSlice(defines)
		wList.Each(
			func(i int, wd *WordDefSt) {
				wd.Def = defines[i]
				answer, _ := arraylist.New(defines...).Find(
					func(_ int, value NLStringSt) bool {
						return value == section.Words.get(wd.Word)
					},
				)
				wd.Answer = NLStringSt{fmt.Sprintf("%c", 'A'+answer)}
			},
		)
		w.Add(wList)
	}
}

func (w *WordMatchSt) ToQuestions() error {
	w.AllQuestions = QuestionSetSt{List: arraylist.New[*QuestionsSt]()}
	words := NLStringListSt{List: arraylist.New[NLStringSt](w.AllWords.Keys()...)}

	w.AllWords.Each(func(word, def NLStringSt) {
		choice := linkedhashset.New[string](word.string)
		for choice.Size() < 4 {
			choice.Add(words.getRandom())
		}

		w.AllQuestions.Add(
			&QuestionsSt{
				Question: NLStringSt{string: "Definition: " + def.string},
				Answer:   1,
				Choices:  WordsSt{List: arraylist.New(choice.Values()...)},
			},
		)
	})
	return nil
}

func (p *PassageCompletionSt) Init(section JSONSectionSt, numTest uint) {
	p.SectionHeadSt = getSectionHead(section)
	p.WordList = StringListSt{arraylist.New[NLStringListSt]()}
	p.Answers = StringListSt{arraylist.New[NLStringListSt]()}

	for i := 0; i < int(numTest); i++ {
		words := arraylist.New[NLStringSt]()
		for _, w := range shuffleSlice(section.WordList.Values()) {
			words.Add(NLStringSt{w})
		}

		p.WordList.Insert(i, NLStringListSt{List: words})

		answers := NLStringListSt{List: arraylist.New[NLStringSt]()}
		p.WordList.get(i).Each(
			func(wi int, _ NLStringSt) {
				answer, _ := p.WordList.get(i).Find(
					func(index int, value NLStringSt) bool {
						return value.string == section.WordList.Values()[wi]
					},
				)
				answers.Insert(wi, NLStringSt{fmt.Sprintf("%c", 'A'+answer)})
			},
		)
		p.Answers.Insert(i, answers)
	}
}

func (c *CompQuestionsSt) Init(section JSONSectionSt, numTest uint) {
	c.SectionHeadSt = getSectionHead(section)
	c.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	c.AllQuestions = section.Questions
}

func (c *CustomSt) Init(section JSONSectionSt, numTest uint) {
	c.SectionHeadSt = getSectionHead(section)
	c.Questions = section.Questions.GetQuestions(numTest, section.NumQuest,
		section.NumCol, section.KeepOrder)
	c.AllQuestions = section.Questions
	c.Answers = section.Answers.List
	c.AnswerText = section.AnswerText.List
}

func (questions QuestionSetSt) GetQuestions(numTest, numQuest, numCol uint,
	keepOrder bool) QuestionListSt {
	newQuestions := QuestionListSt{
		List: arraylist.New[QuestionSetSt](),
	}
	if numQuest > uint(questions.Size()) {
		numQuest = uint(questions.Size())
	}
	for newQuestions.Size() < int(numTest) {
		newQuests := QuestionSetSt{
			List: arraylist.New[*QuestionsSt](),
		}
		questions.QuestionSet(numQuest, numTest, keepOrder).Each(
			func(_ int, q *QuestionsSt) {
				newQuest := &QuestionsSt{
					NumCol:   ternary(q.NumCol != 0, q.NumCol, numCol),
					Question: q.Question,
					Parts:    q.Parts.fixMissing(),
					Required: q.Required,
					Choices:  q.Choices.fixMissing(),
					Answers:  q.Answers.fixMissing(),
				}
				if newQuest.Choices.Size() != 0 {
					answer, _ := newQuest.Choices.Get(int(q.Answer) - 1)
					newQuest.Choices.List = arraylist.New(shuffleSlice(newQuest.Choices.Values())...)
					newQuest.Choices.Each(
						func(ci int, c string) {
							if answer == c {
								newQuest.Answers.List = arraylist.New(fmt.Sprintf("%c", 'A'+ci), c)
							}
						},
					)
				}
				newQuests.Add(newQuest)
			},
		)
		newQuestions.Add(newQuests)
	}

	return newQuestions
}

func (q QuestionSetSt) QuestionSet(numQuest, numTest uint,
	keepOrder bool) QuestionSetSt {
	if keepOrder {
		return QuestionSetSt{List: arraylist.New(q.Values()[:numQuest]...)}
	}

	newSet := arraylist.New(
		arraylist.New(q.Values()...).Select(
			func(_ int, q *QuestionsSt) bool {
				return q.Required
			},
		).Values()...)

	for newSet.Size() < int(numQuest) {
		minUse := numTest
		arraylist.New(q.Values()...).Each(
			func(_ int, q *QuestionsSt) {
				minUse = genfuncs.Min(minUse, q.Used)
			},
		)

		var winner *QuestionsSt
		for found := true; found; found = newSet.Contains(winner) {
			winner = shuffleSlice(arraylist.New(q.Values()...).Select(
				func(_ int, q *QuestionsSt) bool {
					return q.Used == minUse
				},
			).Values())[0]
		}
		newSet.Add(winner)
		winner.Used++
	}

	return QuestionSetSt{List: arraylist.New(shuffleSlice(newSet.Values())...)}
}

func (w WordDistMapSt) WordSet(numQuest, numTest uint) NLStringListSt {
	newSet := arraylist.New[NLStringSt]()
	for newSet.Size() < int(numQuest) {
		minUse := numTest
		w.Each(
			func(word NLStringSt, _ uint) {
				use, _ := w.Get(word)
				minUse = genfuncs.Min(minUse, use)
			},
		)

		var winner NLStringSt
		for found := true; found; found = newSet.Contains(winner) {
			winner = shuffleSlice(w.Select(
				func(w NLStringSt, use uint) bool {
					return use == minUse
				},
			).Keys())[0]
		}
		w.Put(winner, minUse+1)
		newSet.Add(winner)
	}

	return NLStringListSt{List: arraylist.New(shuffleSlice(newSet.Values())...)}
}

func ProcessInclude(section JSONSectionSt, assetdir string) {
	if section.Type == "word-match" {
		ProcessWordsInclude(section, assetdir)
		return
	}

	section.Include.Each(
		func(_ int, inc string) {
			filePath := assetdir + "/" + inc
			_, bytes, err := jsonc.ReadFromFile(filePath)

			if err != nil {
				fmt.Printf("Unable to load include file %s\n", filePath)
				return
			}

			incJSON := make([]*QuestionsSt, 0)
			if err := json.Unmarshal(bytes, &incJSON); err != nil {
				fmt.Printf("ProcessInclude Unable to parse include file %s: %v\n",
					filePath, err)
				return
			}

			section.Questions.Add(incJSON...)
		},
	)

	ProcessIncludeQuestgen(section, assetdir)
	ProcessIncludeAiken(section, assetdir)
}

func ProcessIncludeQuestgen(section JSONSectionSt, assetdir string) {
	section.IncludeQuestgen.Each(
		func(_ int, inc string) {
			filePath := assetdir + "/" + inc
			_, bytes, err := jsonc.ReadFromFile(filePath)

			if err != nil {
				fmt.Printf("Unable to load Questgen include file %s\n", filePath)
				return
			}

			incQuestgenJSON := make([]*QuestgenQuestionSt, 0)
			if err := json.Unmarshal(bytes, &incQuestgenJSON); err != nil {
				fmt.Printf("ProcessInclude Unable to parse Questgen include file %s: %v\n",
					filePath, err)
				return
			}

			for _, q := range incQuestgenJSON {
				newQuest := QuestionsSt{
					Choices: q.Choices,
					Answers: WordsSt{List: arraylist.New[string]()},
				}
				newQuest.Question.string = q.Question
				newQuest.Answers.Add(q.Answer)
				section.Questions.Add(&newQuest)
			}
		},
	)
}

func ProcessIncludeAiken(section JSONSectionSt, assetdir string) {
	section.IncludeAiken.Each(
		func(_ int, inc string) {
			filePath := assetdir + "/" + inc
			incAiken, err := aiken.ReadAiken(filePath)
			if err != nil {
				fmt.Printf("Unable to load Aiken include file %s, %v\n", filePath, err)
				return
			}

			for _, q := range incAiken {
				newQuest := QuestionsSt{
					Choices: WordsSt{List: arraylist.New[string]()},
				}
				newQuest.Question.string = q.Question
				for _, o := range q.Options {
					newQuest.Choices.Add(o.Desc)
					if strings.TrimSpace(q.Answer) == o.Answer {
						newQuest.Answers.List = arraylist.New(o.Answer, o.Desc)
						newQuest.Answer = uint(o.Answer[0]-'A') + 1
					}
				}
				section.Questions.Add(&newQuest)
			}
		},
	)
}

func ProcessWordsInclude(section JSONSectionSt, assetdir string) {
	section.Include.Each(
		func(_ int, inc string) {
			filePath := assetdir + "/" + inc
			_, bytes, err := jsonc.ReadFromFile(filePath)

			if err != nil {
				fmt.Printf("Unable to load include file %s\n", filePath)
				return
			}

			incJSON := new(WordDefMapSt)
			if err := incJSON.UnmarshalJSON(bytes); err != nil {
				fmt.Printf("ProcessWordsInclude Unable to parse include file %s: %v\n",
					filePath, err)
				return
			}

			incJSON.Each(
				func(key, value NLStringSt) {
					section.Words.Put(key, value)
				},
			)
		},
	)
}

func getSectionHead(section JSONSectionSt) SectionHeadSt {
	return SectionHeadSt{
		Type:             section.Type,
		SectionTitle:     section.SectionTitle,
		Title:            section.Title,
		NumLines:         section.NumLines,
		Points:           section.Points,
		NumCol:           section.NumCol,
		NumQuest:         section.NumQuest,
		AnswerLines:      section.AnswerLines,
		QuizBox:          section.QuizBox,
		KeepOrder:        section.KeepOrder,
		Instructions:     section.Instructions.string,
		FormInstructions: section.FormInstructions.string,
		Text:             section.Text.string,
		ColumnHead:       section.ColumnHead,
		Quiz:             false,
	}
}

func (r *ReadingCompSt) GetHead() *SectionHeadSt {
	return &r.SectionHeadSt
}

func (m *MultipleChoiceSt) GetHead() *SectionHeadSt {
	return &m.SectionHeadSt
}

func (w *WordProblemSt) GetHead() *SectionHeadSt {
	return &w.SectionHeadSt
}

func (q *QuizSt) GetHead() *SectionHeadSt {
	return &q.SectionHeadSt
}

func (c *CompQuestionsSt) GetHead() *SectionHeadSt {
	return &c.SectionHeadSt
}

func (w *WordMatchSt) GetHead() *SectionHeadSt {
	return &w.SectionHeadSt
}

func (p *PassageCompletionSt) GetHead() *SectionHeadSt {
	return &p.SectionHeadSt
}

func (c *CustomSt) GetHead() *SectionHeadSt {
	return &c.SectionHeadSt
}
