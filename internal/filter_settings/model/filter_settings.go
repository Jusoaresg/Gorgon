package model

type FilterSettings struct {
	DefaultFilterProfileID *int64 `json:"default_filter_profile_id"`
	UseAliases             bool   `json:"use_aliases"`
	OnlyLatin              bool   `json:"only_latin"`
}

func DefaultFilterSettings() FilterSettings {
	return FilterSettings{
		DefaultFilterProfileID: nil,
		UseAliases:             true,
		OnlyLatin:              true,
	}
}
