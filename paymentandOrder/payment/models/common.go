package models

type CommonModel struct {
	Id           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	LastModified int64  `json:"last_modified" gorm:"index"`
}
