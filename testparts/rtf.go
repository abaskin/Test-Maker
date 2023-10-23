package testparts

import (
	"fmt"

	rtfdoc "github.com/abaskin/Test-Maker/v2/rtf-doc"

	"github.com/chonla/roman-number-go"
)

type RTFDoc struct {
	*rtfdoc.Document
}

const tableWidth = 10000

const charStrFmt = "%c. %s"

func (doc *RTFDoc) Init() {
	doc.Document = rtfdoc.NewDocument()
	doc.SetOrientation(rtfdoc.OrientationPortrait)
	doc.SetFormat(rtfdoc.FormatA4)
}

func (doc *RTFDoc) PageHeader(head *TestHeadSt) {
	// not implemented
}

func (doc *RTFDoc) PageFooter(head *TestHeadSt) {
	t := doc.MakeTable().
		SetWidth(tableWidth).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetBorderColor(rtfdoc.ColorWhite)

	cWidth := t.GetTableCellWidthByRatio(1, 1, 1)
	tr := t.AddTableRow()

	tr.AddDataCell(cWidth[0]).
		AddParagraph().
		SetAlign(rtfdoc.AlignLeft).
		AddText(
			fmt.Sprintf("Gr. %s %s", head.Grade, head.Subject),
			14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)

	tr.AddDataCell(cWidth[1]).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(head.School, 14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)

	tr.AddDataCell(cWidth[2]).
		AddParagraph().
		SetAlign(rtfdoc.AlignRight).
		AddText("Page \\chpgn", 14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)

	doc.AddPageFooter(rtfdoc.FooterAll, t)
}

func (doc *RTFDoc) TestHeader(head *TestHeadSt, sections []SectionSt) {
	t := doc.AddTable().
		SetWidth(tableWidth).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetBorderColor(rtfdoc.ColorWhite)

	tr := t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(head.School,
			14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(fmt.Sprintf("Grade %s, %s, %s", head.Grade, head.Subject, head.RTFTitle),
			14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(fmt.Sprintf("Time Allowed: %d minutes", head.Time),
			12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(fmt.Sprintf("Total Score: %d", head.Points),
			12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		AddText("+", 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorWhite)

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText("Test Sections",
			12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	cWidth := t.GetTableCellWidthByRatio(1, 1)
	for si, s := range sections {
		tr = t.AddTableRow()
		tr.AddDataCell(cWidth[0]).
			AddParagraph().
			SetAlign(rtfdoc.AlignLeft).
			AddText(fmt.Sprintf("%s. %s", roman.NewRoman().ToRoman(si+1), s.GetHead().SectionTitle),
				12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()
		tr.AddDataCell(cWidth[1]).
			AddParagraph().
			SetAlign(rtfdoc.AlignRight).
			AddText(fmt.Sprintf("(%d Points)", s.GetHead().Points),
				12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()
	}

	tr = t.AddTableRow()
	tr.AddDataCell(tableWidth).
		AddParagraph().
		AddText("+", 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorWhite)
}

func (doc *RTFDoc) Sections(student uint, sections []SectionSt, qNum *QuestNumSt) {
	for i, section := range sections {
		doc.sectionHeader(section.GetHead(), i+1)
		qNum.NewSection()
		section.TestRTF(doc, student, qNum)
	}
}

func (doc *RTFDoc) sectionHeader(head *SectionHeadSt, num int) {
	t := doc.AddTable().
		SetWidth(tableWidth).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetBorderColor(rtfdoc.ColorWhite)
	cWidth := t.GetTableCellWidthByRatio(1, 1)

	tr := t.AddTableRow()
	tr.AddDataCell(cWidth[0]).
		AddParagraph().
		SetAlign(rtfdoc.AlignLeft).
		AddText(fmt.Sprintf("Section %s. %s",
			roman.NewRoman().ToRoman(num), head.SectionTitle),
			14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()
	tr.AddDataCell(cWidth[1]).
		AddParagraph().
		SetAlign(rtfdoc.AlignRight).
		AddText(fmt.Sprintf("(%d points)", head.Points),
			14, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	doc.AddParagraph().
		SetAlign(rtfdoc.AlignLeft).
		AddNewLine().
		AddText(head.Instructions, 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
}

func (m *MultipleChoiceSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	questions, _ := m.Questions.Get(int(student))
	rtfQuestions(doc, questions, &m.SectionHeadSt, qNum)
}

func (r *ReadingCompSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	rtfText(doc, r.Title, r.Text)
	questions, _ := r.Questions.Get(int(student))
	rtfQuestions(doc, questions, &r.SectionHeadSt, qNum)
}

func (w *WordProblemSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	questions, _ := w.Questions.Get(int(student))
	rtfQuestions(doc, questions, &w.SectionHeadSt, qNum)
}

func (q *QuizSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	// not implimented
}

func (w *WordMatchSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	t := doc.AddTable().SetWidth(tableWidth)
	t.SetMarginLeft(50).SetMarginRight(50).SetMarginTop(50).SetMarginBottom(50)
	t.SetBorderColor(rtfdoc.ColorWhite)

	cWidth := t.GetTableCellWidthByRatio(1, 3)

	tr := t.AddTableRow()
	tr.AddDataCell(cWidth[0]).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(w.ColumnHead.Values()[0],
			12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	tr.AddDataCell(cWidth[1]).
		AddParagraph().
		SetAlign(rtfdoc.AlignCenter).
		AddText(w.ColumnHead.Values()[1],
			12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()

	rows, _ := w.Get(int(student))
	rows.Each(
		func(i int, row *WordDefSt) {
			tr := t.AddTableRow()
			tr.AddDataCell(cWidth[0]).
				AddParagraph().
				AddText(fmt.Sprintf("%d. %s", qNum.NextNumber(), row.Word.rtfString()),
					12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)

			tr.AddDataCell(cWidth[1]).
				AddParagraph().
				AddText(fmt.Sprintf(charStrFmt, 'A'+i, row.Def.rtfString()),
					12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
		},
	)
	p := doc.AddParagraph()
	p.AddText(" ", 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorWhite)
}

func (pc *PassageCompletionSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	rtfText(doc, pc.Title, pc.Text)

	t := doc.AddTable().SetWidth(tableWidth)
	t.SetMarginLeft(50).SetMarginRight(50).SetMarginTop(50).SetMarginBottom(50)
	t.SetBorderColor(rtfdoc.ColorWhite)

	wordList, _ := pc.WordList.Get(int(student))
	numCols := ternary(pc.NumCol != 0, pc.NumCol, 5)
	cWidth := int(tableWidth / numCols)
	var tr *rtfdoc.TableRow
	wordList.Each(
		func(i int, word NLStringSt) {
			if i%int(numCols) == 0 {
				tr = t.AddTableRow()
			}
			tr.AddDataCell(cWidth).
				AddParagraph().
				AddText(word.string, 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
		},
	)
}

func (c *CompQuestionsSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	questions, _ := c.Questions.Get(int(student))
	rtfQuestions(doc, questions, &c.SectionHeadSt, qNum)
}

func (c *CustomSt) TestRTF(doc *RTFDoc, student uint, qNum *QuestNumSt) {
	questions, _ := c.Questions.Get(0)
	if questions.Size() != 0 {
		questions, _ := c.Questions.Get(int(student))
		rtfQuestions(doc, questions, &c.SectionHeadSt, qNum)
	}
}

func rtfQuestions(doc *RTFDoc, questions QuestionSetSt, head *SectionHeadSt, qNum *QuestNumSt) {
	t := doc.AddTable().
		SetWidth(tableWidth).
		SetMarginLeft(50).
		SetMarginRight(50).
		SetMarginTop(50).
		SetMarginBottom(50).
		SetBorderColor(rtfdoc.ColorWhite)

	questions.Each(
		func(_ int, q *QuestionsSt) {
			t.AddTableRow().
				AddDataCell(tableWidth).
				AddParagraph().
				AddText(fmt.Sprintf("%d. %s", qNum.NextNumber(), q.Question.rtfString()),
					12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)

			cWidth := int(tableWidth / q.NumCol)
			var cRow *rtfdoc.TableRow
			q.Choices.Each(
				func(ci int, c string) {
					if ci%int(q.NumCol) == 0 {
						cRow = t.AddTableRow()
					}
					cRow.AddDataCell(cWidth).
						AddParagraph().
						AddText(fmt.Sprintf(charStrFmt, 'A'+ci, c),
							12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
				},
			)

			q.Parts.Each(
				func(pi int, part NLStringSt) {
					t.AddTableRow().
						AddDataCell(tableWidth).
						AddParagraph().
						AddText(fmt.Sprintf(charStrFmt, 'a'+pi, part.rtfString()),
							12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
				},
			)

			bCol := t.AddTableRow().AddDataCell(tableWidth).AddParagraph()
			bCol.AddText(" ", 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorWhite)
		},
	)
}

func rtfText(doc *RTFDoc, title, text string) {
	p := doc.AddParagraph().SetAlign(rtfdoc.AlignCenter)
	p.AddText(title, 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack).SetBold()
	p.AddNewLine().SetAlign(rtfdoc.AlignLeft)
	p.AddText(text, 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorBlack)
	p.AddNewLine()
	p.AddText("+", 12, rtfdoc.FontTimesNewRoman, rtfdoc.ColorWhite)
	p.AddNewLine()
}
