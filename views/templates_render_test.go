package views

import (
	"bytes"
	"testing"

	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
	filterProfileModel "github.com/jusoaresg/gorgon/internal/filter_profile/model"
	filterSettingsModel "github.com/jusoaresg/gorgon/internal/filter_settings/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showAliasModel "github.com/jusoaresg/gorgon/internal/show_aliases/model"
	showSettingsModel "github.com/jusoaresg/gorgon/internal/show_settings/model"
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
				InfoUrl:   "https://example.com/release/1",
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

func TestRenderFilterSettings(t *testing.T) {
	tmpl := NewTemplate()

	profileID := int64(1)
	data := struct {
		FilterSettings filterSettingsModel.FilterSettings
		Profiles       []filterProfileModel.FilterProfile
	}{
		FilterSettings: filterSettingsModel.FilterSettings{
			DefaultFilterProfileID: &profileID,
			UseAliases:             true,
			OnlyLatin:              true,
		},
		Profiles: []filterProfileModel.FilterProfile{
			{ID: 1, Name: "HD"},
			{ID: 2, Name: "SD"},
		},
	}

	var buf bytes.Buffer
	if err := tmpl.templates.ExecuteTemplate(&buf, "filterSettings", data); err != nil {
		t.Fatalf("failed to render filterSettings: %v", err)
	}

	for _, want := range []string{
		"Default Filter Profile",
		"HD",
		"SD",
		"Save Profile",
		"only_latin",
		"json-enc-custom",
		"Placeholder reference",
		"{season:00}",
		"{absolute}",
	} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("rendered output missing %q", want)
		}
	}
}

func TestRenderEditShowModal(t *testing.T) {
	tmpl := NewTemplate()

	profileID := int64(1)
	data := struct {
		Show     showModel.Show
		Profiles []filterProfileModel.FilterProfile
		Settings showSettingsModel.ShowSettings
		Aliases  []showAliasModel.ShowAlias
	}{
		Show: showModel.Show{
			ID:   42,
			Name: "Dragon Ball",
		},
		Profiles: []filterProfileModel.FilterProfile{
			{ID: 1, Name: "HD"},
			{ID: 2, Name: "SD"},
		},
		Settings: showSettingsModel.ShowSettings{
			FilterProfileID: &profileID,
			UseAliases:      true,
			OnlyLatin:       true,
		},
		Aliases: []showAliasModel.ShowAlias{
			{ID: 1, Alias: "DBZ", Source: "user"},
			{ID: 2, Alias: "ドラゴンボール", Source: "tvmaze"},
		},
	}

	var buf bytes.Buffer
	if err := tmpl.templates.ExecuteTemplate(&buf, "edit-show-modal", data); err != nil {
		t.Fatalf("failed to render edit-show-modal: %v", err)
	}

	for _, want := range []string{
		"Edit Series",
		"Dragon Ball",
		"/api/v1/database/show-settings/42",
		"filter_profile_id",
		"use_aliases",
		"only_latin",
		"DBZ",
		"Delete",
		"ドラゴンボール",
		"synced",
		"Custom Aliases",
		"Delete Series",
		"Save Changes",
	} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("rendered output missing %q", want)
		}
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
