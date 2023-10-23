package testparts

import (
	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/daichi-m/go18ds/sets/hashset"
)

func (w *WordMatchSt) TestForm(gf *GoogleFormSt, student uint) error {
	words := w.getWords(student)
	points := w.Points / uint(words.Size())
	questions := QuestionSetSt{List: arraylist.New[*QuestionsSt]()}

	wordDef, _ := w.Get(int(student))
	wordDef.Each(
		func(_ int, wDef *WordDefSt) {
			choice := hashset.New(wDef.Word.string)
			for choice.Size() < 5 {
				choice.Add(words.getRandom())
			}

			def, _ := w.getDefs(student).Get(int(wDef.Answer.string[0] - 'A'))

			questions.Add(
				&QuestionsSt{
					Question: def,
					Answers:  WordsSt{List: arraylist.New("", wDef.Word.string)},
					Choices:  WordsSt{List: arraylist.New(choice.Values()...)},
				},
			)
		},
	)

	if err := gf.AddSection(w.SectionTitle, w.FormInstructions); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (c *CompQuestionsSt) TestForm(gf *GoogleFormSt, student uint) error {
	questions, _ := c.Questions.Get(int(student))
	points := c.Points / uint(questions.Size())
	if err := gf.AddSection(c.SectionTitle, c.Text); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (w *WordProblemSt) TestForm(gf *GoogleFormSt, student uint) error {
	questions, _ := w.Questions.Get(int(student))
	points := w.Points / uint(questions.Size())
	if err := gf.AddSection(w.SectionTitle, w.Text); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (m *MultipleChoiceSt) TestForm(gf *GoogleFormSt, student uint) error {
	questions, _ := m.Questions.Get(int(student))
	points := m.Points / uint(questions.Size())
	if err := gf.AddSection(m.SectionTitle, m.FormInstructions); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (r *ReadingCompSt) TestForm(gf *GoogleFormSt, student uint) error {
	questions, _ := r.Questions.Get(int(student))
	points := r.Points / uint(questions.Size())
	if err := gf.AddSection(r.SectionTitle, r.Text); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (q *QuizSt) TestForm(gf *GoogleFormSt, student uint) error {
	questions, _ := q.Questions.Get(int(student))
	points := q.Points / uint(questions.Size())
	if err := gf.AddSection(q.SectionTitle, q.Text); err != nil {
		return err
	}
	return gf.AddQuestions(questions, int64(points))
}

func (c *CustomSt) TestForm(gf *GoogleFormSt, student uint) error {
	return nil
}

func (p *PassageCompletionSt) TestForm(gf *GoogleFormSt, student uint) error {
	return nil
}
