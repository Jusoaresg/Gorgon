package schema

type TorrentContent struct {
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	Size         int64   `json:"size"`
	Progress     float64 `json:"progress"`
	Priority     int     `json:"priority"`
	Is_Seed      bool    `json:"is_seed"`
	Piece_range  []int64 `json:"piece_range"`
	Availability float64 `json:"availability"`
}
