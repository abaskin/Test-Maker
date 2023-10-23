package testparts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/nwillc/genfuncs"
	"github.com/rwestlund/gotex"
	"gorm.io/gorm"
	"muzzammil.xyz/jsonc"
)

func GetTestJson(fileName, assetdir string) (TestJSONSt, error) {
	filePath := assetdir + "/" + fileName
	_, bytes, err := jsonc.ReadFromFile(filePath)

	if err != nil {
		log.Println("Unable to load test file, error: ", err)
		return TestJSONSt{}, err
	}

	testJSON := TestJSONSt{}

	if err := json.Unmarshal(bytes, &testJSON); err != nil {
		log.Println("Unable to parse test file, error: ", err)
		return TestJSONSt{}, err
	}

	testJSON.Logo, _ = filepath.Abs(assetdir + "/" + testJSON.Logo)

	return testJSON, nil
}

func MakeSections(testJSON TestJSONSt, assetDir string, showAll bool) ([]SectionSt, error) {
	sections := make([]SectionSt, 0)

	for _, section := range testJSON.Sections {
		log.Println("Making section:", section.SectionTitle)
		section.AnswerText.fixMissing()
		section.WordList.fixMissing()
		section.Include.fixMissing()
		section.IncludeQuestgen.fixMissing()
		section.IncludeAiken.fixMissing()
		section.Answers.fixMissing()
		section.ColumnHead.fixMissing()
		section.Words.fixMissing()
		section.Questions.fixMissing()
		section.AssetDir = assetDir
		numTests := uint(testJSON.Students.Size())
		ProcessInclude(section, assetDir)

		if section.NumQuest == 0 {
			section.NumQuest = uint(genfuncs.Min(section.Questions.Size(), section.Words.Size()))
		}

		if section.NumCol == 0 {
			section.NumCol = 4
		}

		section.Questions.Each(
			func(_ int, q *QuestionsSt) {
				if q.Required {
					q.Used = numTests + 1
				}
			},
		)

		if showAll {
			numTests = 1
			section.NumQuest = uint(section.Questions.Size())
		}

		switch section.Type {
		case "word-match":
			if showAll {
				section.NumQuest = uint(section.Words.Size())
			}
			w := new(WordMatchSt)
			w.Init(section, numTests)
			sections = append(sections, w)

		case "multiple-choice":
			m := new(MultipleChoiceSt)
			m.Init(section, numTests)
			sections = append(sections, m)

		case "word-problem":
			w := new(WordProblemSt)
			w.Init(section, numTests)
			sections = append(sections, w)

		case "quiz":
			q := new(QuizSt)
			q.Init(section, numTests)
			sections = append(sections, q)

		case "reading-comprehension":
			r := new(ReadingCompSt)
			r.Init(section, numTests)
			sections = append(sections, r)

		case "comprehension-questions":
			c := new(CompQuestionsSt)
			c.Init(section, numTests)
			sections = append(sections, c)

		case "passage-completion":
			p := new(PassageCompletionSt)
			p.Init(section, numTests)
			sections = append(sections, p)

		case "custom":
			c := new(CustomSt)
			c.Init(section, numTests)
			sections = append(sections, c)

		default:
			log.Printf("Error: unknown section type %s\n", section.Type)
		}
	}

	return sections, nil
}

func MakeBundle(dsn string, testJSON TestJSONSt, sections []SectionSt) (*arraylist.List[*TestBundleSt], error) {
	testHead := TestHeadSt{
		Subject:  testJSON.Subject,
		Grade:    testJSON.Grade,
		School:   testJSON.School,
		Title:    testJSON.Title,
		RTFTitle: testJSON.RTFTitle,
		Logo:     testJSON.Logo,
		Time:     testJSON.Time,
		Date:     testJSON.Date,
		NoKey:    testJSON.NoKey,
		Quiz:     false,
		Points:   0,
		Dsn:      dsn,
	}

	for _, s := range sections {
		if s.GetHead().Quiz {
			testHead.Quiz = true
		}
		testHead.Points += s.GetHead().Points
	}

	testBundle := arraylist.New[*TestBundleSt]()
	testJSON.Students.Each(
		func(si int, student string) {
			testBundle.Add(
				&TestBundleSt{
					Student:    student,
					StudentNum: uint(si),
					TestHeadSt: &testHead,
				},
			)
		},
	)

	return testBundle, nil
}

func (bundle *TestBundleSt) Create(dsn string, pathStrings PathStrSt,
	flags FlagsSt, test TestSt) error {
	if bundle.Quiz {
		if err := bundle.createQuiz(pathStrings, flags, test); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	if err := bundle.dbImport(dsn, pathStrings, flags, test); err != nil {
		log.Println(err)
	}

	if err := bundle.createRTF(pathStrings, flags, test); err != nil {
		log.Println(err)
	}

	if err := bundle.createPDF(pathStrings, flags, test); err != nil {
		log.Println(err)
	}

	if err := bundle.createForm(pathStrings, flags, test); err != nil {
		log.Println(err)
	}

	if !flags.ShowAll {
		if err := bundle.createAnswerSheet(pathStrings, flags, test, false); err != nil {
			log.Println(err)
		}
	}

	if !bundle.NoKey {
		if err := bundle.createAnswerSheet(pathStrings, flags, test, true); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (distroTest *TestBundleSt) CreateDistro(pathStrings PathStrSt, flags FlagsSt, test TestSt) {
	if flags.CreateDistro {
		outStr := []string{string(test.Template), ""}
		outStr = append(outStr, DocumentBegin(distroTest)...)
		outStr = append(outStr, QuestDistoSheet(test.TestJSON, test.Sections)...)
		outStr = append(outStr, DocumentEnd(distroTest)...)

		if err := makePDF(pathStrings.Workdir, pathStrings.Outdir, "question", "distro",
			strings.Join(outStr, "\n"), flags.SaveTex); err != nil {
			log.Println("Unable to create distro file, error: ", err)
		}
	}
}

func (bundle *TestBundleSt) createQuiz(pathStrings PathStrSt, flags FlagsSt, test TestSt) error {
	if flags.CreatePDF {
		testID := strings.ReplaceAll(bundle.Student, " ", "")
		qNum := MakeQuestNum(!flags.ContinuousNumbering)
		outStr := []string{string(test.Template), ""}
		outStr = append(outStr, DocumentBegin(bundle)...)
		outStr = append(outStr, QuizSheet(bundle, test.Sections, &qNum)...)
		outStr = append(outStr, DocumentEnd(bundle)...)

		if err := makePDF(pathStrings.Workdir, pathStrings.Outdir, testID, "quiz",
			strings.Join(outStr, "\n"), flags.SaveTex); err != nil {
			return fmt.Errorf("unable to create quiz file, error: %w", err)
		}
	}
	return nil
}

func (bundle *TestBundleSt) dbImport(dsn string, pathStrings PathStrSt,
	flags FlagsSt, test TestSt) error {
	if flags.DBImport {
		log.Println("Connecting to database")
		db, err := OpenCockroachDB(dsn, false)
		if err != nil {
			return fmt.Errorf("unable connect to database, error: %w", err)
		}
		defer func() {
			CloseCockroachDB(db)
			log.Println("Disconnected from database")
		}()
		log.Println("Connected to database")

		newTest := GormTest{
			Title:        test.TestJSON.Title,
			Length:       test.TestJSON.Time,
			MinQuestions: test.TestJSON.MinQuestions,
			Questions:    make([]GormQuestion, 0),
		}

		for _, section := range test.Sections {
			points := section.GetHead().Points / section.GetHead().NumQuest
			allQuestions := QuestionSetSt{}

			switch section.GetHead().Type {
			case "word-match":
				wordMatch := section.(*WordMatchSt)
				wordMatch.ToQuestions()
				allQuestions = wordMatch.AllQuestions

			case "multiple-choice":
				allQuestions = section.(*MultipleChoiceSt).AllQuestions

			// case "word-problem":
			// case "quiz":
			// case "reading-comprehension":
			// case "comprehension-questions":
			// case "passage-completion":
			// case "custom":

			default:
				log.Printf("Warning: unsupported section type %s\n", section.GetHead().Type)
				allQuestions.Size() // make the linter happy
				continue
			}

			log.Printf("%d questions to add", allQuestions.Size())
			allQuestions.Each(func(_ int, question *QuestionsSt) {
				choices := make([]GormQuestionChoice, 0)
				question.Choices.Each(func(index int, choice string) {
					choices = append(choices,
						GormQuestionChoice{
							Choice: choice,
							Answer: index+1 == int(question.Answer),
						})
				})
				newTest.Questions = append(newTest.Questions,
					GormQuestion{
						Required: question.Required,
						Question: question.Question.CleanString(),
						Points:   points,
						Choices:  choices,
					})
			})
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			// do some database operations in the transaction
			// (use 'tx' from this point, not 'db')
			class := GormClass{}
			if err := tx.
				Where(GormClass{Subject: test.TestJSON.Subject}).
				First(&class).Error; err != nil {
				return err
			}

			class.Tests = append(class.Tests, newTest)

			if err := tx.Save(&class).Error; err != nil {
				// return any error will rollback
				return err
			}
			log.Printf("inserted test record for %s %s\n", class.Subject, test.TestJSON.Title)

			// return nil will commit the whole transaction
			return nil
		}); err != nil {
			return fmt.Errorf("database error: %w", err)
		}
	}
	return nil
}

func (bundle *TestBundleSt) createRTF(pathStrings PathStrSt, flags FlagsSt, test TestSt) error {
	if flags.CreateRtf {
		testID := strings.ReplaceAll(bundle.Student, " ", "")
		qNum := MakeQuestNum(!flags.ContinuousNumbering)
		rtf := new(RTFDoc)
		rtf.Init()
		rtf.TestHeader(bundle.TestHeadSt, test.Sections)
		rtf.Sections(bundle.StudentNum, test.Sections, &qNum)
		rtf.PageFooter(bundle.TestHeadSt)
		return makeRTF(pathStrings.Outdir, testID, "test", rtf)
	}
	return nil
}

func (bundle *TestBundleSt) createPDF(pathStrings PathStrSt, flags FlagsSt, test TestSt) error {
	if flags.CreatePDF {
		testID := strings.ReplaceAll(bundle.Student, " ", "")
		qNum := MakeQuestNum(!flags.ContinuousNumbering)
		outStr := []string{string(test.Template), ""}
		outStr = append(outStr, DocumentBegin(bundle)...)
		outStr = append(outStr, TestSheet(bundle, test.Sections, &qNum)...)
		outStr = append(outStr, DocumentEnd(bundle)...)

		if err := makePDF(pathStrings.Workdir, pathStrings.Outdir, testID, "test",
			strings.Join(outStr, "\n"), flags.SaveTex); err != nil {
			return fmt.Errorf("unable to create test file, error: %w", err)
		}
	}
	return nil
}

func (bundle *TestBundleSt) createForm(pathStrings PathStrSt, flags FlagsSt, test TestSt) error {
	if flags.CreateForm {
		title := bundle.RTFTitle
		if title == "" {
			title = bundle.Title
		}
		form, err := NewGoogleForm(bundle.Dsn).Create(title, bundle.Student, "")
		if err != nil {
			return fmt.Errorf("unable to create Google Form, error: %w", err)
		}
		for _, section := range test.Sections {
			section.TestForm(form, bundle.StudentNum)
			if err != nil {
				log.Printf("Google Form section, error: %s\n", err.Error())
			}
		}
	}
	return nil
}

func (bundle *TestBundleSt) createAnswerSheet(pathStrings PathStrSt, flags FlagsSt, test TestSt, isKey bool) error {
	if flags.CreatePDF {
		testID := strings.ReplaceAll(bundle.Student, " ", "")
		qNum := MakeQuestNum(!flags.ContinuousNumbering)
		outStr := []string{string(test.Template), ""}
		outStr = append(outStr, DocumentBegin(bundle)...)
		outStr = append(outStr, AnswerSheet(bundle, test.Sections, isKey, false, &qNum)...)
		outStr = append(outStr, DocumentEnd(bundle)...)

		sheetType := "answer"
		if isKey {
			sheetType = "key"
		}

		if err := makePDF(pathStrings.Workdir, pathStrings.Outdir, testID, sheetType,
			strings.Join(outStr, "\n"), flags.SaveTex); err != nil {
			return fmt.Errorf("unable to create answer sheet file, , key: %t, error: %w", isKey, err)
		}
	}
	return nil
}

func makePDF(workdir, outdir, testID, otype, testTex string, saveTex bool) error {
	if saveTex {
		testPath := fmt.Sprintf("%s/%s-%s.tex", workdir, testID, otype)

		tFile, err := os.Create(testPath)
		if err != nil {
			log.Println("Unable to create test file, error: ", err)
			return err
		}
		testwriter := bufio.NewWriter(tFile)
		testwriter.WriteString(testTex)
		testwriter.Flush()
		tFile.Close()
	}

	testPdf, err := gotex.Render(testTex,
		gotex.Options{
			Command:   "/Library/TeX/texbin/pdflatex",
			Runs:      0,
			Texinputs: "",
		})

	if err != nil {
		log.Println("render failed ", err)
		return err
	}

	testOut := fmt.Sprintf("%s/%s-%s.pdf", outdir, testID, otype)
	os.WriteFile(testOut, testPdf, 0644)

	return nil
}

func makeRTF(outdir, testID, otype string, doc *RTFDoc) error {
	testPath := fmt.Sprintf("%s/%s-%s.rtf", outdir, testID, otype)

	tFile, err := os.Create(testPath)
	if err != nil {
		log.Println("Unable to create RTF test file, error: ", err)
		return err
	}
	testwriter := bufio.NewWriter(tFile)
	testwriter.WriteString(string(doc.Export()))
	testwriter.Flush()
	tFile.Close()

	return nil
}
