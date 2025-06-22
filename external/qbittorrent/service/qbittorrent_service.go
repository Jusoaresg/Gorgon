package service

import (
	"bytes"
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/internal/episode_content/model"
	"github.com/jusoaresg/gorgon/pkg/services"
	"log/slog"
	"mime/multipart"
	"net/url"
)

type QBittorrentService struct {
	APIService     *services.APIService
	Logger         *slog.Logger
	username       string
	password       string
	host           string
	port           string
	DownloadFolder string // The client has to have access to this folder
	sid            string // API KEY
}

func NewQBittorrentService(logger *slog.Logger) (*QBittorrentService, error) {
	configFile, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	host := configFile.QBittorrentHost
	port := configFile.QBittorrentPort
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
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			q.sid = cookie.Value
			return nil
		}
	}
	return fmt.Errorf("QBittorrent SID not found")
}

func (q *QBittorrentService) SidVerification() error {
	if q.sid == "" {
		req := schema.QBittorrentLoginRequest{
			Username: q.username,
			Password: q.password,
		}
		if err := q.Login(&req); err != nil {
			return fmt.Errorf("error while logging on qbittorrent: %w", err)
		}

	}
	return nil
}

func (q *QBittorrentService) Logout() error {
	return nil
}

func (q *QBittorrentService) AddTorrent(url string) error {

	if err := q.SidVerification(); err != nil {
		return err
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	writer.WriteField("urls", url)
	writer.WriteField("savepath", q.DownloadFolder)
	writer.WriteField("category", "gorgon")
	writer.WriteField("skip_checking", "false")
	writer.WriteField("root_folder", "false")
	writer.WriteField("paused", "true")
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Cookie":       fmt.Sprintf("SID=%s", q.sid),
	}

	_, err := q.APIService.Post("/api/v2/torrents/add", requestBody.Bytes(), nil, headers)
	if err != nil {
		return err
	}

	return nil
}

// Filter can be "all", "downloading", "seeding", "completed", paused", "active", "inactive", "resumed", "staled", "stalled_uploading", stalled_downloading", "errored"
func (q *QBittorrentService) CheckTorrents(filter string, response *[]schema.CheckTorrentResponse) error {

	if err := q.SidVerification(); err != nil {
		return err
	}

	headers := map[string]string{
		"Cookie": fmt.Sprintf("SID=%s", q.sid),
	}

	url := fmt.Sprintf("/api/v2/torrents/info?filter=%s", filter)
	//TODO: Logger message with url and maybe the headers

	if err := q.APIService.GetWithHeaders(url, &response, headers); err != nil {
		return fmt.Errorf("error while get torrent info: %w", err)
	}

	return nil
}

func (q *QBittorrentService) CheckTorrentsWithHash(filter, hash string, response *[]schema.CheckTorrentResponse) error {

	if err := q.SidVerification(); err != nil {
		return err
	}

	headers := map[string]string{
		"Cookie": fmt.Sprintf("SID=%s", q.sid),
	}

	encodedHash := url.QueryEscape(hash)
	url := fmt.Sprintf("/api/v2/torrents/info?filter=%s&hashes=%s", filter, encodedHash)
	//TODO: Logger message with url and maybe the headers

	if err := q.APIService.GetWithHeaders(url, &response, headers); err != nil {
		return fmt.Errorf("error while get torrent info: %w", err)
	}

	return nil
}

func (q *QBittorrentService) GetContent(hash string) ([]model.EpisodeContent, error) {
	if err := q.SidVerification(); err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Cookie": fmt.Sprintf("SID=%s", q.sid),
	}

	url := fmt.Sprintf("/api/v2/torrents/files?hash=%s", hash)
	var content []model.EpisodeContent
	if err := q.APIService.GetWithHeaders(url, &content, headers); err != nil {
		return nil, err
	}

	return content, nil
}
