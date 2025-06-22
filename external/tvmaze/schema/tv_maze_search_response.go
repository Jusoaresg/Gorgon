package schema

import "github.com/jusoaresg/gorgon/pkg/schemas/dtos"

type TvMazeResponse struct {
	Score float64      `json:"score"`
	Show  dtos.ShowDto `json:"show"`
}
