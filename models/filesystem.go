package models

import (
	"time"
)

type Filesystem struct {
	ID        int    `json:"id" gorm:"primarykey"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Filesystem) TableName() string {
	return "filesystem"
}
