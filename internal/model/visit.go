package model

type Visit struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	LinkID    uint64 `gorm:"index"`
	IP        string `gorm:"size:64;not null"`
	UserAgent string `gorm:"size:255;not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

func (Visit) TableName() string {
	return "visits"
}
