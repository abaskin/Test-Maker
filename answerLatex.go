package testparts

import (
	"fmt"
	"strings"
)

func (r *ReadingCompSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(r.SectionTitle, r.Points, "")
	questions := r.Questions.get(student)
	quest, _ := questions.Get(0)
	if quest.Choices.Size() == 0 {
		return append(outStr, answerLines(questions, isKey, r.GetHead().NumLines)...)
	}
	start := qNum.CurrentNumber() + 1
	end := qNum.CurrentNumber() + uint32(r.NumQuest)
	qNum.AddNumber(uint32(r.NumQuest))
	outStr = append(outStr, fmt.Sprintf(`\answerBox{%d}{%d}`, start, end))
	if isKey {
		outStr = append(outStr,
			fmt.Sprintf(`{%s%s}`, strings.Repeat("x", int(start)-1),
				questionAnswerString(questions)))
	}
	return outStr
}

func (m *MultipleChoiceSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	start := qNum.CurrentNumber() + 1
	end := qNum.CurrentNumber() + uint32(m.NumQuest)
	qNum.AddNumber(uint32(m.NumQuest))
	outStr := []string{
		strings.Join(testSectionBegin(m.SectionTitle, m.Points, ""), "\n"),
		fmt.Sprintf(`\answerBox{%d}{%d}`, start, end),
	}
	if isKey {
		outStr = append(outStr,
			fmt.Sprintf(`{%s%s}`, strings.Repeat("x", int(start)-1),
				questionAnswerString(m.Questions.get(student))))
	}
	return outStr
}

func (w *WordProblemSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(w.SectionTitle, w.Points, "")

	if showAll {
		return append(outStr, []string{
			`\begin{enumerate}`,
			strings.Join(stringsUsing(w.Questions.get(student).Values(), func(value *QuestionsSt) string {
				return fmt.Sprintf(`\item %s`, strings.Join(value.Answers.Values(), ", "))
			}), "\n"),
			`\end{enumerate}`,
		}...)
	}

	return append(outStr, answerTable(w.Questions.get(student), isKey)...)
}

func (q *QuizSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(q.SectionTitle, q.Points, "")

	if showAll {
		return append(outStr, []string{
			`\begin{enumerate}`,
			strings.Join(stringsUsing(q.Questions.get(student).Values(), func(value *QuestionsSt) string {
				return fmt.Sprintf(`\item %s`, strings.Join(value.Answers.Values(), ", "))
			}), "\n"),
			`\end{enumerate}`,
		}...)
	}

	return append(outStr, answerTable(q.Questions.get(student), isKey)...)
}

func (w *WordMatchSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	start := qNum.CurrentNumber() + 1
	end := qNum.CurrentNumber() + uint32(w.NumQuest)
	qNum.AddNumber(uint32(w.NumQuest))
	outStr := testSectionBegin(w.SectionTitle, w.Points, "")
	outStr = append(outStr, fmt.Sprintf(`\answerBox{%d}{%d}`, start, end))
	if isKey {
		outStr = append(outStr, fmt.Sprintf(`{%s%s}`,
			strings.Repeat("x", int(start)-1),
			strings.Join(w.getAnswers(student).Values(), "")))
	}
	return outStr
}

func (p *PassageCompletionSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(p.SectionTitle, p.Points, "")
	wordList, _ := p.WordList.Get(int(student))
	start := qNum.CurrentNumber() + 1
	end := qNum.CurrentNumber() + uint32(wordList.Size())
	qNum.AddNumber(uint32(wordList.Size()))
	outStr = append(outStr, fmt.Sprintf(`\answerBox{%d}{%d}`, start, end))
	if isKey {
		answers, _ := p.Answers.Get(int(student))
		outStr = append(outStr, fmt.Sprintf(`{%s%s}`,
			strings.Repeat("x", int(start)-1), answers.join("")))
	}
	return outStr
}

func (c *CompQuestionsSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(c.SectionTitle, c.Points, "")

	return append(outStr, answerLines(c.Questions.get(student), isKey, c.GetHead().NumLines)...)
}

func (c *CustomSt) AnswerLatex(isKey, showAll bool, student uint, qNum *QuestNumSt) []string {
	outStr := testSectionBegin(c.SectionTitle, c.Points, "")

	if isKey {
		return append(outStr, []string{
			`\large`,
			strings.Join(c.Answers.Values(), "\n"),
			`\normalsize`,
		}...)
	}

	return append(outStr, c.AnswerText.Values()...)
}
