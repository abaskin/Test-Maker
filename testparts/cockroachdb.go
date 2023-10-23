package testparts

import (
	"log"
	"os"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func OpenCockroachDB(dsn string, autoMigrate bool) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Millisecond * 500, // Slow SQL threshold
			LogLevel:      logger.Info,            // Log level
		},
	)
	db, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		return nil, err
	}

	// only need to run this when the schema changes
	if autoMigrate {
		err = db.AutoMigrate(&GormTest{}, &GormQuestion{}, &GormQuestionChoice{},
			GromPreference{}, GormClass{}, &GormStudent{}, GormTestAttempt{},
			GormTestSession{})
	}

	return db, err
}

func CloseCockroachDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err == nil {
		err = sqlDB.Close()
	}
	return err
}

type Preference struct {
	dsn string
}

func newPreference(dsn string) *Preference {
	return &Preference{
		dsn: dsn,
	}
}

func (p *Preference) Set(key, value string) error {
	db, err := OpenCockroachDB(p.dsn, false)
	if err != nil {
		return err
	}
	defer CloseCockroachDB(db)

	return db.Create(
		&GromPreference{
			Key:   key,
			Value: value,
		}).Error
}

func (p *Preference) Get(key string) (string, error) {
	db, err := OpenCockroachDB(p.dsn, false)
	if err != nil {
		return "", err
	}
	defer CloseCockroachDB(db)

	pref := GromPreference{}
	result := db.Where(&GromPreference{Key: key}).First(&pref)
	return pref.Value, result.Error
}

func (p *Preference) Del(key string) error {
	db, err := OpenCockroachDB(p.dsn, false)
	if err != nil {
		return err
	}
	defer CloseCockroachDB(db)

	return db.Delete(&GromPreference{Key: key}).Error
}
