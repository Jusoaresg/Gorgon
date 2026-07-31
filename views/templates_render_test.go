package views

import (
	"bytes"
	"testing"

	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
)

type renderData struct {
	Items        []service.DownloadItem
	ErrorMessage string
}

func TestRenderDownloadsPage(t *testing.T) {
	tmpl := NewTemplate()

	items := []service.DownloadItem{
		{
			Torrent: schema.CheckTorrentResponse{
				Name:    "My.Show.S01E01.1080p.WEB-DL",
				Hash:    "abc123",
				State:   "downloading",
				Category: "gorgon",
				Progress: 0.42,
				Size:     734003200,
				Completed: 308281344,
				DlSpeed:  5242880,
				UpSpeed:  1048576,
				NumSeeds: 12,
				NumLeechs: 3,
				Eta:      98,
			},
			Episode: &service.EpisodeInfo{
				EpisodeID: 1,
				ShowID:    1,
				Name:      "Pilot",
				Season:    1,
				Number:    1,
				ShowName:  "Breaking Test",
				ShowImage: "https://example.com/img.jpg",
			},
		},
		{
			Torrent: schema.CheckTorrentResponse{
				Name:     "Some.Other.Torrent",
				Hash:     "xyz",
				State:    "stalledDL",
				Category: "gorgon",
				Progress: 0.0,
				Eta:      -1,
			},
		},
		{
			Torrent: schema.CheckTorrentResponse{
				Name:     "My.Show.S01E02.1080p.WEB-DL",
				Hash:     "donehash",
				State:    "uploading",
				Category: "gorgon",
				Progress: 1.0,
				Size:     500000000,
				Completed: 500000000,
			},
			Episode: &service.EpisodeInfo{
				EpisodeID: 2,
				ShowID:    1,
				Name:      "Second Episode",
				Season:    1,
				Number:    2,
				Tracking:  "snatched",
				ShowName:  "Breaking Test",
			},
		},
	}

	var buf bytes.Buffer
	if err := tmpl.templates.ExecuteTemplate(&buf, "download-items", PageData{Data: renderData{Items: items}}); err != nil {
		t.Fatalf("failed to render download-items: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"Breaking Test",
		"S01E01",
		"Pilot",
		"42%",
		"Downloading",
		"5.0 MB/s",
		"1.0 MB/s",
		"12 / 3",
		"1m 38s",
		"294.0 MB / 700.0 MB",
		"Some.Other.Torrent",
		"Stalled",
		"Waiting to Import",
		"100%",
		"hx-trigger=\"every 2s\"",
	} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("rendered output missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderDownloadsEmpty(t *testing.T) {
	tmpl := NewTemplate()

	var buf bytes.Buffer
	if err := tmpl.templates.ExecuteTemplate(&buf, "download-items", PageData{Data: renderData{}}); err != nil {
		t.Fatalf("failed to render empty download-items: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("No active downloads")) {
		t.Errorf("expected empty state, got:\n%s", buf.String())
	}
}

func TestRenderDownloadsFullLayout(t *testing.T) {
	tmpl := NewTemplate()

	var buf bytes.Buffer
	err := tmpl.templates.ExecuteTemplate(&buf, "layout", PageData{
		TemplateName: "downloads",
		Styles:       []string{"downloads.css"},
		Data: renderData{
			Items: []service.DownloadItem{
				{
					Torrent: schema.CheckTorrentResponse{
						Name:    "My.Show.S01E01.1080p.WEB-DL",
						Hash:    "abc123",
						State:   "downloading",
						Category: "gorgon",
						Progress: 0.42,
						Size:     734003200,
						Completed: 308281344,
						DlSpeed:  5242880,
						NumSeeds: 12,
						NumLeechs: 3,
						Eta:      98,
					},
					Episode: &service.EpisodeInfo{
						EpisodeID: 1,
						ShowID:    1,
						Name:      "Pilot",
						Season:    1,
						Number:    1,
						ShowName:  "Breaking Test",
						ShowImage: "https://example.com/img.jpg",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to render full layout: %v", err)
	}

	for _, want := range []string{
		"href=\"/downloads\"",
		"class=\"sidebar-icon\"",
		"downloads.css",
		"Breaking Test",
		"S01E01",
		"42%",
		"hx-trigger=\"every 2s\"",
	} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("rendered layout missing %q", want)
		}
	}
}
