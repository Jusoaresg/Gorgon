package schema

// FilterSettingsDto is the request payload for updating global filter
// settings. Pointer fields distinguish "not provided" (nil) from an explicit
// value, so partial updates merge cleanly with the persisted settings. The
// default profile id arrives as a string: an empty string clears the default,
// a numeric string sets it.
type FilterSettingsDto struct {
	DefaultFilterProfileID *string `json:"default_filter_profile_id"`
	UseAliases             *bool   `json:"use_aliases"`
	OnlyLatin              *bool   `json:"only_latin"`
}
