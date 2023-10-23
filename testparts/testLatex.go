package testparts

import (
	"fmt"
	"strings"

	"github.com/daichi-m/go18ds/lists/arraylist"
)

func (m *MultipleChoiceSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(m.SectionTitle, m.Points, m.Instructions)
	questions, _ := m.Questions.Get(int(student))
	qNum.NewSection()
	return append(outStr, questionsLatex(questions, *m.GetHead(), false, qNum)...)
}

func (r *ReadingCompSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(r.SectionTitle, r.Points, r.Instructions)
	questions, _ := r.Questions.Get(int(student))
	return append(outStr, []string{
		fmt.Sprintf(`\qtitle{%s}`, r.Title),
		r.Text,
		`\\`,
		strings.Join(questionsLatex(questions, *r.GetHead(), false, qNum), "\n"),
	}...)
}

func (w *WordProblemSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(w.SectionTitle, w.Points, w.Instructions)
	outStr = append(outStr, w.Text)
	questions, _ := w.Questions.Get(int(student))
	return append(outStr, questionsLatex(questions, *w.GetHead(), false, qNum)...)
}

func (q *QuizSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := make([]string, 0)
	if q.QuizBox {
		questions, _ := q.Questions.Get(int(student))
		outStr = append(outStr, []string{
			`\begin{quizgrid}`,
			strings.Join(stringsUsing(questions.Values(), func(qz *QuestionsSt) string {
				return strings.Join([]string{
					`\quizbox{`, qz.Question.string, `}`}, "\n")
			}), "\n"),
			`\end{quizgrid}`}...)
		return outStr
	}
	questions, _ := q.Questions.Get(int(student))
	return append(outStr, questionsLatex(questions, *q.GetHead(), true, qNum)...)
}

func (w *WordMatchSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(w.SectionTitle, w.Points, w.Instructions)
	outStr = append(outStr, w.Text)
	words := w.getWords(student)
	defines := w.getDefs(student)

	enumi := fmt.Sprintf(`\setcounter{enumi}{%d}`, qNum.CurrentNumber())
	qNum.AddNumber(uint32(words.Size()))
	return append(outStr, []string{
		fmt.Sprintf(`\qwm{%s}{%s} {`, w.ColumnHead.Values()[0], w.ColumnHead.Values()[1]),
		`\qwmWord`,
		enumi,
		optionList(`\qwmItem {%s}`, words),
		`\qwmColEnd } { \qwmDef`,
		optionList(`\qwmItem {%s}`, defines),
		`\qwmColEnd }`,
	}...)
}

func (p *PassageCompletionSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(p.SectionTitle, p.Points, p.Instructions)
	outStr = append(outStr, []string{
		`\qformat{\hfill} \begin{questions}`,
		`\titledquestion{} \fullwidth{`,
		fmt.Sprintf(`\qtitle{%s}`, p.Title),
		p.Text,
		`}`,
	}...)

	wordList, _ := p.WordList.Get(int(student))
	return append(outStr, []string{
		fmt.Sprintf(`\begin{qchoices}(%d)`, ternary(p.NumCol != 0, p.NumCol, 5)),
		optionList(`\choice %s`, wordList),
		`\end{qchoices} \end{questions} \noqformat`,
	}...)
}

func (c *CompQuestionsSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(c.SectionTitle, c.Points, c.Instructions)
	questions, _ := c.Questions.Get(int(student))
	return append(outStr, questionsLatex(questions, *c.GetHead(), false, qNum)...)
}

func (c *CustomSt) TestLatex(student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(c.SectionTitle, c.Points, c.Instructions)
	outStr = append(outStr, c.Text)
	questions, _ := c.Questions.Get(0)
	if questions.Size() != 0 {
		questions, _ := c.Questions.Get(int(student))
		outStr = append(outStr, questionsLatex(questions, *c.GetHead(), false, qNum)...)
	}
	return outStr
}

func optionList(format string, options NLStringListSt) string {
	optionList := arraylist.New[string]()
	options.Each(
		func(_ int, value NLStringSt) {
			optionList.Add(fmt.Sprintf(format, value.string))
		},
	)
	return strings.Join(optionList.Values(), "\n")
}
