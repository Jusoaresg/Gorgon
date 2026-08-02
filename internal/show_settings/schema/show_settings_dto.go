package schema

import (
	showSettingsModel "github.com/jusoaresg/gorgon/internal/show_settings/model"
)

type ShowSettingsDto struct {
	FilterProfileID *int64 `json:"filter_profile_id"`
	UseAliases      bool   `json:"use_aliases"`
	OnlyLatin       bool   `json:"only_latin"`
}

func ToShowSettingsDto(settings showSettingsModel.ShowSettings) ShowSettingsDto {
	return ShowSettingsDto{
		FilterProfileID: settings.FilterProfileID,
		UseAliases:      settings.UseAliases,
		OnlyLatin:       settings.OnlyLatin,
	}
}
