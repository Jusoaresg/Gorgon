package schema

import (
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"
)

type SearchResult struct {
	Show    dtos.ShowDto
	IsAdded bool
}
