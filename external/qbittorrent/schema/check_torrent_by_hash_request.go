package schema

type CheckTorrentByHashRequest struct {
	Status string `json:"status"`
	Hash   string `json:"hash"`
}
