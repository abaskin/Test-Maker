package testparts

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/nwillc/genfuncs"
)

func testQR(test *TestBundleSt) (string, string) {
	md5 := md5.Sum([]byte(test.Title))
	return strings.Join([]string{
		test.Student,
		fmt.Sprintf("%02x%02x%02x%02x", md5[3], md5[2], md5[1], md5[0]),
	}, `-('_')-`), test.Student
}

func DocumentBegin(test *TestBundleSt) []string {
	qrText, studentName := testQR(test)
	return []string{
		`\begin{document}`,
		`\testSetFooter`,
		fmt.Sprintf(`{%s}{%s}{%s}`, test.Grade, test.Subject, test.School),
		`{Page \thepage}`,
		fmt.Sprintf(`\testSetHeader {%s}{%s}`, qrText, studentName),
	}
}

func DocumentEnd(test *TestBundleSt) []string {
	return []string{
		`\end{document}`,
	}
}

func QuizSheet(quiz *TestBundleSt, sections []SectionSt, qNum *QuestNumSt) []string {
	outStr := []string{
		`\begin{quiz}`,
		fmt.Sprintf(`{%s}{%s}{%s}{%d}`, quiz.Student, quiz.Title, quiz.Date, quiz.Points),
		`\newline`,
	}

	for _, s := range sections {
		outStr = append(outStr, s.TestLatex(quiz.StudentNum, qNum)...)
	}

	return append(outStr, `\end{quiz}`)
}

func TestSheet(test *TestBundleSt, sections []SectionSt, qNum *QuestNumSt) []string {
	outStr := []string{
		`\begin{test}`,
		fmt.Sprintf(`{%s}{0.75}`, test.Logo),
		fmt.Sprintf(`{%s}{%s}{%s}`, test.Grade, test.Title, test.Subject),
		fmt.Sprintf(`{%d}{%d}`, test.Time, test.Points),
		`{ \begin{answerSections}`,
		strings.Join(stringsUsing(sections, func(value SectionSt) string {
			h := value.GetHead()
			return fmt.Sprintf(`\answerSectionLine {%s}{%d}`, h.SectionTitle, h.Points)
		}), "\n"),
		`\end{answerSections} }`,
	}

	for _, s := range sections {
		qNum.NewSection()
		outStr = append(outStr, s.TestLatex(test.StudentNum, qNum)...)
	}

	return append(outStr, `\end{test}`)
}

func AnswerSheet(test *TestBundleSt, sections []SectionSt, isKey, showAll bool,
	qNum *QuestNumSt) []string {
	outStr := make([]string, 0)

	switch {
	case isKey:
		outStr = append(outStr, `\begin{answerKey}`)
	default:
		outStr = append(outStr, []string{
			`\begin{answerSheet}`,
			fmt.Sprintf(`{%s}{0.75}`, test.Logo),
			fmt.Sprintf(`{%s}{%s}{%s}`, test.Grade, test.Title, test.Subject),
			fmt.Sprintf(`{%d}{%d}{%s}`, test.Points, test.Time, test.Student),
			fmt.Sprintf(`{%s}`, test.Date),
		}...)
	}

	for _, s := range sections {
		qNum.NewSection()
		outStr = append(outStr, s.AnswerLatex(isKey, showAll, test.StudentNum, qNum)...)
	}

	outStr = append(outStr, "")

	return append(outStr, ternary(isKey, `\end{answerKey}`, `\end{answerSheet}`))
}

func QuestDistoSheet(test TestJSONSt, sections []SectionSt) []string {
	outStr := []string{}

	for _, s := range sections {
		graph := s.DistribLatex(test.Title, uint(test.Students.Size()))
		if len(graph) > 0 {
			outStr = append(outStr, graph...)
		}
	}

	return outStr
}

func questionsLatex(questions QuestionSetSt, head SectionHeadSt, isQuiz bool, qNum *QuestNumSt) []string {
	outStr := []string{`\begin{questions}`}
	outStr = append(outStr, fmt.Sprintf(`\setcounter{question}{%d}`, qNum.CurrentNumber()))
	qNum.AddNumber(uint32(questions.Size()))

	questions.Each(
		func(_ int, q *QuestionsSt) {
			outStr = q.Begin(outStr)

			numCol := uint(4)
			if head.NumCol != 0 {
				numCol = head.NumCol
			}
			if q.NumCol != 0 {
				numCol = q.NumCol
			}

			if q.Choices.Size() != 0 {
				outStr = q.Choice(outStr, numCol)
			}

			if q.Parts.Size() != 0 {
				outStr = q.Part(outStr)
			}

			outStr = append(outStr, `\end{minipage}`)
			outStr = q.Lines(outStr, isQuiz, head)
			outStr = append(outStr, " ")
		},
	)

	return append(outStr, `\end{questions}`)
}

func answerTable(questions QuestionSetSt, showAnswers bool) []string {
	numCol := 1
	questions.Each(
		func(_ int, q *QuestionsSt) {
			if q.Answers.Size() == 0 && showAnswers {
				return
			}
			numCol = genfuncs.Min(numCol, q.Parts.Size())
		},
	)

	outStr := []string{
		`\noindent`,
		fmt.Sprintf(
			`\begin{tabularx}{\textwidth}{@{\rule[-5.25mm]{0pt}{12mm}}|Y|*{%d}{C|}}`,
			numCol),
		`\hline`,
	}

	header := sequenceUsing(
		[]string{},
		func(value int) string {
			return fmt.Sprintf("%c", 'a'+value)
		},
		0, numCol)

	header = append([]string{`Question`}, header...)
	outStr = append(outStr, strings.Join(header, " & ")+` \\ \hline`)

	questions.Each(func(qi int, q *QuestionsSt) {
		aLine := sequenceUsing(
			[]string{},
			func(value int) string {
				return " "
			},
			0, numCol)

		if showAnswers {
			copy(aLine, q.Answers.Values())
		}
		aLine = append([]string{fmt.Sprintf("%d", qi+1)}, aLine...)
		outStr = append(outStr, strings.Join(aLine, " & ")+` \\ \hline`)
	})

	return append(outStr, `\end{tabularx}`)
}

func testSectionBegin(title string, points uint, inst string) []string {
	return []string{
		"",
		fmt.Sprintf(`\testSection{%s}`, title),
		fmt.Sprintf(`{%d}{%s}`, points, inst),
	}
}

func questionAnswerString(questions QuestionSetSt) string {
	return strings.Join(stringsUsing(questions.Values(), func(q *QuestionsSt) string {
		a, _ := q.Answers.Get(0)
		return a
	}), "")
}

func answerLines(questions QuestionSetSt, isKey bool, numLines string) []string {
	if isKey {
		return []string{
			`\begin{enumerate}`,
			`\large`,
			strings.Join(stringsUsing(questions.Values(), func(value *QuestionsSt) string {
				return fmt.Sprintf(`\item %s`, strings.Join(value.Answers.Values(), "\n"))
			}), "\n"),
			`\normalsize`,
			`\end{enumerate}`,
		}
	}

	lines := ternary(numLines == "", "3.5cm", numLines)
	return []string{fmt.Sprintf(`\answerLines{%d}{%s}`, questions.Size(), lines)}
}

func (q *QuestionsSt) Begin(outStr []string) []string {
	return append(outStr, []string{
		`\begin{minipage}{\linewidth}`,
		`\question`,
		q.Question.string,
	}...)
}

func (q *QuestionsSt) Choice(outStr []string, numCol uint) []string {
	return append(outStr, []string{
		fmt.Sprintf(`\begin{qchoices}(%d)`, numCol),
		strings.Join(stringsUsing(q.Choices.Values(), func(value string) string {
			return fmt.Sprintf(`\choice %s`, value)
		}), "\n"),
		`\end{qchoices}`,
	}...)
}

func (q *QuestionsSt) Part(outStr []string) []string {
	return append(outStr, []string{
		`\begin{parts}`,
		strings.Join(stringsUsing(q.Parts.values(),
			func(value string) string {
				return fmt.Sprintf(`\part %s`, value)
			},
		), "\n"),
		`\end{parts}`,
	}...)
}

func (q *QuestionsSt) Lines(outStr []string, isQuiz bool, head SectionHeadSt) []string {
	if isQuiz {
		if head.NumLines == "" {
			switch {
			case head.AnswerLines:
				return append(outStr, `\fillwithlines{\stretch{1}}`)
			default:
				return append(outStr, `\vspace*{\stretch{1}}`)
			}
		}

		switch {
		case head.AnswerLines:
			return append(outStr, fmt.Sprintf(`\fillwithlines{%s}`, head.NumLines))
		default:
			return append(outStr, fmt.Sprintf(`\vspace*{%s}}`, head.NumLines))
		}
	}
	return append(outStr, `\vspace{0.25cm}`)
}
