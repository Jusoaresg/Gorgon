package model

type Anime struct {
	Id                int `gorm:";primaryKey"`
	Aid               int `gorm:"uniqueIndex"`
	EpisodeCount      int
	Description       string
	AiringTime        string
	NextAiringEpisode *Episode `gorm:"foreignKey:Aid"`
	Title             Title    `gorm:"foreignKey:Aid"`
	Genres            string
	BannerImage       string          `json:"bannerImage"`
	CoverImage        string          `json:"coverImage"`
	Relations         RelationWrapper `gorm:"foreignKey:Aid"`
	Status            string

	InstalledEps []Episode `gorm:"foreignKey:Aid"`
}

type RelationWrapper struct {
	Id    int       `gorm:"primaryKey"`
	Aid   int       `gorm:"index"`
	Edges []Related `gorm:"foreignKey:RelationsId"`
}

type Related struct {
	Id          int   `gorm:"primaryKey"`
	RelationsId int   `gorm:"index"`
	Aid         int   `gorm:"index"`
	Title       Title `gorm:"foreignKey:Aid"`
	Episodes    int
	Type        string
	Format      string
	Status      string
}

type Title struct {
	Id      int `gorm:"primaryKey"`
	Aid     int `gorm:"index"`
	Romaji  string
	English string
}

type Episode struct {
	Id       int `gorm:"primaryKey"`
	Aid      int `gorm:"index"`
	Title    string
	Episode  int
	AiringAt int
	Aired    bool `gorm:"not null"`
}

// OTHER MODEL
type Season struct {
	Id            int `gorm:"primaryKey"`
	Aid           int `gorm:"index"`
	MainAnimeAid  int `gorm:"index"`
	SeasonAnimeId int `gorm:"uniqueIndex"`
	SeasonNumber  int `gorm:"not null"`
}
