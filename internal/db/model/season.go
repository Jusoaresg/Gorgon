package model

type Season struct {
	ID       uint `gorm:"primaryKey"`
	ShowId   int  `gorm:"index"`
	SeasonId int  `gorm:"index"`
	Number   int
}
