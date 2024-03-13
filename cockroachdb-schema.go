package testparts

import (
	"math/rand"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ----- Preferences -----

type GromPreference struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex"`
	Value string
}

// ----- Classes -----

type GormClass struct {
	gorm.Model
	Subject  string `gorm:"uniqueIndex:index:idx_class"`
	Sections []GormClassSection
	Tests    []GormTest
}

type GormClassSection struct {
	gorm.Model
	Section     string `gorm:"uniqueIndex:index:idx_section"`
	GormClassID uint   `gorm:"uniqueIndex:index:idx_section"`
	Students    []GormStudent
}

type GormStudent struct {
	gorm.Model
	FamilyName         string
	GivenName          string
	GormClassSectionID uint
}

// ----- Tests -----

type GormTest struct {
	gorm.Model
	Title        string `gorm:"index:idx_test"`
	Length       uint
	MinQuestions uint
	GormClassID  uint
	Sessions     []GormTestSession
	Attempts     []GormTestAttempt
	Questions    []GormQuestion
}

type GormQuestion struct {
	gorm.Model
	Required   bool
	Question   string
	Points     uint
	GormTestID uint
	Choices    []GormQuestionChoice
}

type GormQuestionChoice struct {
	gorm.Model
	GormQuestionID uint
	Choice         string
	Feedback       string
	Answer         bool
}

type GormTestAttempt struct {
	gorm.Model
	GormStudentID uint
	GormTestID    uint
	Score         float64
	AttemptStart  time.Time
	AttemptEnd    time.Time
	Answers       datatypes.JSON
}

type GormTestSession struct {
	gorm.Model
	GormTestID         uint
	GormClassSectionID uint
	QuestionTime       time.Duration
	StartDateTime      time.Time
	EndDateTime        time.Time
}

func (q *GormQuestion) ShuffleChoices() {
	rand.Shuffle(len(q.Choices),
		func(i, j int) {
			q.Choices[i], q.Choices[j] = q.Choices[j], q.Choices[i]
		})
}

func (t *GormTest) ShuffleQuestions() {
	rand.Shuffle(len(t.Questions),
		func(i, j int) {
			t.Questions[i], t.Questions[j] = t.Questions[j], t.Questions[i]
		})
}
