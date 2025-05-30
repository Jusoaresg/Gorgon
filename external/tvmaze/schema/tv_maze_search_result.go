package schema

import (
	"gorgon/pkg/schemas/dtos"
)

type SearchResult struct {
	Show    dtos.ShowDto
	IsAdded bool
}
