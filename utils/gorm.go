package utils

import "gorm.io/gorm"

func Transaction(db *gorm.DB, f func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
