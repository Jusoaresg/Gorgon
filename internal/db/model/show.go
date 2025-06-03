package model

type Show struct {
	ID        int    `gorm:"primaryKey"`
	ShowID    int    `gorm:"index"`
	Name      string `gorm:"uniqueIndex"`
	Type      string
	Language  string
	Status    string
	Premiered string
	Ended     string
	Rating    *float64
	Summary   string
	Updated   int

	Schedule Schedule  `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`
	Seasons  []Season  `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`
	Episodes []Episode `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`

	Externals Externals `gorm:"embedded"`
	Image     Image     `gorm:"embedded"`

	Genres string `gorm:"type:text"` // Serializar e desserializar manualmente
}

type Schedule struct {
	ID     uint `gorm:"primaryKey"`
	ShowId int  `gorm:"index"`
	Time   string
	Days   string `gorm:"type:text"`
}

type Externals struct {
	Tvrage   int
	Thetvdvb int
	Imdb     string
}

type Image struct {
	Medium   string
	Original string
}
