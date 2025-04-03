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

	// Relacionamentos
	Schedule Schedule  `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`
	Seasons  []Season  `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`
	Episodes []Episode `gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE"`

	// Structs aninhadas
	Externals Externals `gorm:"embedded"`
	Image     Image     `gorm:"embedded"`

	// Campos que precisam ser serializados
	Genres string `gorm:"type:text"` // Serializar e desserializar manualmente
}

type Season struct {
	ID       uint `gorm:"primaryKey"`
	ShowId   int  `gorm:"index"`
	SeasonId int  `gorm:"index"`
	Number   int
}

type Episode struct {
	ID       uint `gorm:"primaryKey"`
	ShowId   int  `gorm:"index"`
	Name     string
	Summary  string
	Type     string
	Number   int
	Season   int
	AirDate  string
	AirStamp string
	AirTime  string
}

type Schedule struct {
	ID     uint `gorm:"primaryKey"`
	ShowId int  `gorm:"index"`
	Time   string
	Days   string `gorm:"type:text"` // Serializar e desserializar manualmente
}

// Structs aninhadas
type Externals struct {
	Tvrage   int
	Thetvdvb int
	Imdb     string
}

type Image struct {
	Medium   string
	Original string
}
