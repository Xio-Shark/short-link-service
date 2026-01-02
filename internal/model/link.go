package model

type Link struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"size:16;uniqueIndex"`
	OriginalURL string `gorm:"size:2048;not null"`
	ExpireAt    int64  `gorm:"not null;default:0"`
	Status      int    `gorm:"not null;default:1"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

func (Link) TableName() string {
	return "links"
}
