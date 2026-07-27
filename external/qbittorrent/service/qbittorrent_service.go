package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/internal/episode_content/model"
	"github.com/jusoaresg/gorgon/pkg/services"
)

type ErrQBittorrentHostPortNotSet interface {
	Error() string
	IsProwlarrConfigWarn() bool
}
type QBittorrentHostPortNotSet struct{}

func (m *QBittorrentHostPortNotSet) Error() string {
	return "QBittorrent Host or Port not set"
}

func (m *QBittorrentHostPortNotSet) IsProwlarrConfigWarn() bool {
	return true
}

type addTorrentResponse struct {
	AddedTorrentIds []string `json:"added_torrent_ids"`
	FailureCount    int      `json:"failure_count"`
	PendingCount    int      `json:"pending_count"`
	SuccessCount    int      `json:"success_count"`
}

type QBittorrentService struct {
	APIService     *services.APIService
	Logger         *slog.Logger
	username       string
	password       string
	host           string
	port           string
	DownloadFolder string
	sid            string
	cookieName     string
	loggedIn       bool
}

func NewQBittorrentService(logger *slog.Logger) (*QBittorrentService, error) {
	configFile, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	host := configFile.QBittorrentHost
	port := configFile.QBittorrentPort

	if host == "" || port == "" {
		return nil, &QBittorrentHostPortNotSet{}
	}

	downloadFolder := configFile.QBittorrentDownloadFolder

	return &QBittorrentService{
		APIService:     services.NewAPIService(fmt.Sprintf("%s:%s", host, port), logger),
		Logger:         logger,
		username:       configFile.QBittorrentUsername,
		password:       configFile.QBittorrentPassword,
		DownloadFolder: downloadFolder,
		host:           host,
		port:           port,
	}, nil
}

func (q *QBittorrentService) Name() string {
	return "qBittorrent"
}

func (q *QBittorrentService) Login(request *schema.QBittorrentLoginRequest) error {
	form := url.Values{}
	form.Add("username", request.Username)
	form.Add("password", request.Password)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      fmt.Sprintf("%s:%s", q.host, q.port),
	}

	resp, err := q.APIService.Post("/api/v2/auth/login", form.Encode(), nil, headers)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid credentials")
	}

	if resp.StatusCode == http.StatusNoContent {
		for _, cookie := range resp.Cookies() {
			if strings.HasPrefix(cookie.Name, "SID") || strings.HasPrefix(cookie.Name, "QBT_SID") {
				q.sid = cookie.Value
				q.cookieName = cookie.Name
				q.loggedIn = true
				return nil
			}
		}
		q.loggedIn = true
		return nil
	}

	for _, cookie := range resp.Cookies() {
		if strings.HasPrefix(cookie.Name, "SID") || strings.HasPrefix(cookie.Name, "QBT_SID") {
			q.sid = cookie.Value
			q.cookieName = cookie.Name
			q.loggedIn = true
			return nil
		}
	}
	return fmt.Errorf("QBittorrent SID not found (status=%d)", resp.StatusCode)
}

func (q *QBittorrentService) SidVerification() error {
	if q.loggedIn {
		return nil
	}
	req := schema.QBittorrentLoginRequest{
		Username: q.username,
		Password: q.password,
	}
	if err := q.Login(&req); err != nil {
		return fmt.Errorf("error while logging on qbittorrent: %w", err)
	}
	return nil
}

func (q *QBittorrentService) cookieHeader() string {
	if q.cookieName != "" {
		return fmt.Sprintf("%s=%s", q.cookieName, q.sid)
	}
	return fmt.Sprintf("SID=%s", q.sid)
}

func (q *QBittorrentService) CheckConnection() error {
	if err := q.SidVerification(); err != nil {
		return err
	}
	return nil
}

func (q *QBittorrentService) Logout() error {
	return nil
}

func (q *QBittorrentService) AddTorrent(torrentUrl string) error {
	if err := q.SidVerification(); err != nil {
		return err
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	writer.WriteField("urls", torrentUrl)
	writer.WriteField("category", "gorgon")
	writer.WriteField("skip_checking", "false")
	writer.WriteField("root_folder", "false")
	writer.WriteField("paused", "true")
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Cookie":       q.cookieHeader(),
	}

	var addResp addTorrentResponse
	_, err := q.APIService.Post("/api/v2/torrents/add", requestBody.Bytes(), &addResp, headers)
	if err != nil {
		return err
	}

	if addResp.FailureCount > 0 {
		return fmt.Errorf("failed to add torrent (failures=%d)", addResp.FailureCount)
	}

	return nil
}

func (q *QBittorrentService) DeleteTorrent(hash string, deleteFile bool) error {
	if err := q.SidVerification(); err != nil {
		return err
	}

	form := url.Values{}
	form.Add("hashes", hash)
	form.Add("deleteFiles", fmt.Sprintf("%t", deleteFile))

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       q.cookieHeader(),
	}

	resp, err := q.APIService.Post("/api/v2/torrents/delete", form.Encode(), nil, headers)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete torrent (status=%d)", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentService) CheckTorrents(filter string, response *[]schema.CheckTorrentResponse) error {
	if err := q.SidVerification(); err != nil {
		return err
	}

	headers := map[string]string{
		"Cookie": q.cookieHeader(),
	}

	endpoint := fmt.Sprintf("/api/v2/torrents/info?filter=%s", filter)
	q.Logger.Debug("Calling QBittorrent API", slog.String("url", endpoint), slog.String("filter", filter))

	if err := q.APIService.GetWithHeaders(endpoint, response, headers); err != nil {
		return fmt.Errorf("error while get torrent info: %w", err)
	}

	return nil
}

func (q *QBittorrentService) CheckTorrentsWithHash(filter, hash string, response *[]schema.CheckTorrentResponse) error {
	if err := q.SidVerification(); err != nil {
		return err
	}

	headers := map[string]string{
		"Cookie": q.cookieHeader(),
	}

	encodedHash := url.QueryEscape(hash)
	endpoint := fmt.Sprintf("/api/v2/torrents/info?filter=%s&hashes=%s", filter, encodedHash)

	if err := q.APIService.GetWithHeaders(endpoint, response, headers); err != nil {
		return fmt.Errorf("error while get torrent info: %w", err)
	}

	return nil
}

func (q *QBittorrentService) GetContent(hash string) ([]model.EpisodeContent, error) {
	if err := q.SidVerification(); err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Cookie": q.cookieHeader(),
	}

	endpoint := fmt.Sprintf("/api/v2/torrents/files?hash=%s", hash)

	resp, err := q.APIService.GetWithHeadersRaw(endpoint, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get torrent files (status=%d)", resp.StatusCode)
	}

	var files []schema.TorrentContent
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("error decoding torrent files: %w", err)
	}

	var content []model.EpisodeContent
	for _, f := range files {
		content = append(content, model.EpisodeContent{
			Name:    f.Name,
			Size:    float64(f.Size),
			Is_Seed: f.Is_Seed,
		})
	}

	return content, nil
}
