package util

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func StartTX(db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	err := tx.Error
	if err != nil {
		logrus.WithError(err).Panic("unable to start transaction")
	}

	return tx
}

func Commit(db *gorm.DB) bool {
	err := db.Commit().Error
	if err != nil {
		logrus.WithError(err).Error("unable to commit transaction")
		return false
	}

	return true
}

func Rollback(db *gorm.DB) {
	err := db.Rollback().Error
	if err != nil {
		logrus.WithError(err).Panic("unable to rollback transaction")
	}
}
