package testparts

import (
	"fmt"
)

func (r *ReadingCompSt) DistribLatex(testTitle string, numTest uint) []string {
	return questDistrib(r.AllQuestions, r.Title, testTitle, r.Points, numTest)
}

func (m *MultipleChoiceSt) DistribLatex(testTitle string, numTest uint) []string {
	return questDistrib(m.AllQuestions, m.Title, testTitle, m.Points, numTest)
}

func (w *WordProblemSt) DistribLatex(testTitle string, numTest uint) []string {
	return questDistrib(w.AllQuestions, w.Title, testTitle, w.Points, numTest)
}

func (q *QuizSt) DistribLatex(testTitle string, numTest uint) []string {
	return questDistrib(q.AllQuestions, q.Title, testTitle, q.Points, numTest)
}

func (c *CompQuestionsSt) DistribLatex(testTitle string, numTest uint) []string {
	return questDistrib(c.AllQuestions, c.Title, testTitle, c.Points, numTest)
}

func (w *WordMatchSt) DistribLatex(testTitle string, numTest uint) []string {
	return []string{}
}

func (p *PassageCompletionSt) DistribLatex(testTitle string, numTest uint) []string {
	return []string{}
}

func (c *CustomSt) DistribLatex(testTitle string, numTest uint) []string {
	return []string{}
}

func questDistrib(questions QuestionSetSt, title, testTitle string,
	points, yMax uint) []string {
	outStr := testSectionBegin(
		ternary(title == "", testTitle, title),
		points, "")

	outStr = append(outStr, []string{
		`\begin{tikzpicture}`,
		fmt.Sprintf(
			`\begin{axis} [ybar ,bar width=15pt, xmin=0, xmax=%d, ymin=0, ymax=%d, width=\linewidth]`,
			questions.Size()+1, yMax),
		`\addplot coordinates {`,
	}...)

	questions.Each(func(qi int, q *QuestionsSt) {
		outStr = append(outStr, fmt.Sprintf(`(%d, %d)`, qi+1, q.Used))
	})

	return append(outStr, []string{
		`};`,
		`\end{axis}`,
		`\end{tikzpicture}`,
	}...)
}
