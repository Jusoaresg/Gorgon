package settings

import (
	"net/http"
	"strings"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/views"
	"github.com/labstack/echo/v4"
)

type SettingType int

const (
	GorgonSettings SettingType = iota
	ProwlarrSettings
	TorrentSettings
)

var SettingTypes = map[string]SettingType{
	"gorgon":   GorgonSettings, // Default Type
	"prowlarr": ProwlarrSettings,
	"torrent":  TorrentSettings,
}

type SettingsData struct {
	Type       SettingType
	TypeString string
	Settings   any
}

func SettingsWithoutParam(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", views.PageData{
		TemplateName: "settings",
		Data: SettingsData{
			Type:       GorgonSettings,
			TypeString: strings.Join([]string{"gorgon", "Settings"}, ""),
		},
		Styles: []string{"settings.css"},
	})
}

func SettingsWithParam(c echo.Context) error {
	settingsType := c.Param("type")

	settingType, ok := SettingTypes[settingsType]
	if !ok {
		settingType = GorgonSettings
		settingsType = "gorgon"
	}

	config, err := config.LoadConfig()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "layout", nil)
	}

	return c.Render(http.StatusOK, "layout", views.PageData{
		TemplateName: "settings",
		Data: SettingsData{
			Type:       settingType,
			TypeString: strings.Join([]string{settingsType, "Settings"}, ""),
			Settings:   config,
		},
		Styles: []string{"settings.css"},
	})
}
