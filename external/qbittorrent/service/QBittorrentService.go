package service

import (
	"bytes"
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/pkg/services"
	"log/slog"
	"mime/multipart"
	"net/url"
)

type QBittorrentService struct {
	APIService *services.APIService
	Logger     *slog.Logger
	username   string
	password   string
	host       string
	port       string
	sid        string // API KEY
}

func NewQBittorrentService(logger *slog.Logger) (*QBittorrentService, error) {
	configFile, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	host := configFile.QBittorrentHost
	port := configFile.QBittorrentPort

	return &QBittorrentService{
		APIService: services.NewAPIService(fmt.Sprintf("%s:%s", host, port), logger),
		Logger:     logger,
		username:   configFile.QBittorrentUsername,
		password:   configFile.QBittorrentPassword,
		host:       host,
		port:       port,
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
			fmt.Println("SID found:", q.sid)
			return nil
		}
	}
	return fmt.Errorf("QBittorrent SID not found")
}

func (q *QBittorrentService) Logout() error {
	return nil
}

func (q *QBittorrentService) AddTorrent(url string) error {

	if q.sid == "" {
		req := schema.QBittorrentLoginRequest{
			Username: q.username,
			Password: q.password,
		}
		if err := q.Login(&req); err != nil {
			return fmt.Errorf("error while logging on qbittorrent: %w", err)
		}

	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	writer.WriteField("urls", url)
	writer.WriteField("savepath", "/downloads")
	writer.WriteField("category", "anime")
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

	if q.sid == "" {
		req := schema.QBittorrentLoginRequest{
			Username: q.username,
			Password: q.password,
		}
		if err := q.Login(&req); err != nil {
			return fmt.Errorf("error while logging on qbittorrent: %w", err)
		}

	}
	headers := map[string]string{
		"Cookie": fmt.Sprintf("SID=%s", q.sid),
	}

	if err := q.APIService.GetWithHeaders("/api/v2/torrents/info", &response, headers); err != nil {
		return fmt.Errorf("error while get torrent info: %w", err)
	}

	return nil
}
