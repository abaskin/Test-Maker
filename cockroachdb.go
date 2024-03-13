package testparts

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

// ---- Preferences ----

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
