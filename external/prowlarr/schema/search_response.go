package schema

import (
	"fmt"
	"time"
)

type SearchResponse struct {
	Title       string  `json:"title"`
	Guid        string  `json:"guid"`
	Age         float32 `json:"age"`
	AgeHours    float32 `json:"ageHours"`
	AgeMinutes  float32 `json:"ageMinutes"`
	Size        int64   `json:"size"`
	Filename    string  `json:"fileName"`
	IndexerId   int     `json:"indexerId"`
	Indexer     string  `json:"indexer"`
	InfoUrl     string  `json:"infoUrl"`
	InfoHash    string  `json:"infoHash"`
	PublishDate string  `json:"publishDate"`
	MagnetUrl   string  `json:"magnetUrl"`
	Seeders     int     `json:"seeders"`
	Leechers    int     `json:"leechers"`
	Protocol    string  `json:"protocol"`
}

func (s *SearchResponse) FormatSize() string {
	size := s.Size
	units := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0

	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}

	return fmt.Sprintf("%v %s", size, units[i])
}

func (s *SearchResponse) FormatAge() string {
	if s.Age >= 1 {
		return fmt.Sprintf("%.1f days", s.Age)
	} else if s.AgeHours >= 1 {
		return fmt.Sprintf("%.1f hours", s.AgeHours)
	}
	return fmt.Sprintf("%.1f minutes", s.AgeMinutes)
}

func (s *SearchResponse) FormatPublishDate() string {
	t, err := time.Parse(time.RFC3339, s.PublishDate)
	if err != nil {
		return "Invalid date"
	}
	return t.Format("02/01/2006 15:04")
}
