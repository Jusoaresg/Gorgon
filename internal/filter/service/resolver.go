package service

import (
	"database/sql"
	"errors"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/filter"
	filterProfileModel "github.com/jusoaresg/gorgon/internal/filter_profile/model"
	filterProfileRepository "github.com/jusoaresg/gorgon/internal/filter_profile/repository"
	filterSettingsRepository "github.com/jusoaresg/gorgon/internal/filter_settings/repository"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
	"github.com/jusoaresg/gorgon/utils"
)

var latinRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_!?',.:]+$`)

// EffectiveSettings is the resolved filtering configuration for a show,
// merging per-show settings over the global defaults. A show_settings row
// fully overrides the global toggles and, when present, the profile.
type EffectiveSettings struct {
	FilterProfileID *int64
	UseAliases      bool
	OnlyLatin       bool
}

func ResolveSettings(db *sqlx.DB, showID int64) (EffectiveSettings, error) {
	settingsRepo := filterSettingsRepository.NewFilterSettingsRepository(db)

	global, err := settingsRepo.Get()
	if err != nil {
		return EffectiveSettings{}, err
	}

	settings := EffectiveSettings{
		FilterProfileID: global.DefaultFilterProfileID,
		UseAliases:      global.UseAliases,
		OnlyLatin:       global.OnlyLatin,
	}

	showSettingsRepo := showSettingsRepository.NewShowSettingsRepository(db)
	showSettings, err := showSettingsRepo.GetByShowID(showID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return settings, nil
		}
		return EffectiveSettings{}, err
	}

	if showSettings.FilterProfileID != nil {
		settings.FilterProfileID = showSettings.FilterProfileID
	}
	settings.UseAliases = showSettings.UseAliases
	settings.OnlyLatin = showSettings.OnlyLatin

	return settings, nil
}

func ResolveProfile(db *sqlx.DB, settings EffectiveSettings) (*filter.Profile, error) {
	if settings.FilterProfileID == nil {
		return nil, nil
	}

	profileRepo := filterProfileRepository.NewFilterProfileRepository(db)
	profileModel, patterns, err := profileRepo.GetByID(*settings.FilterProfileID)
	if err != nil {
		return nil, err
	}

	return ToProfile(profileModel, patterns), nil
}

func ToProfile(p filterProfileModel.FilterProfile, patterns []filterProfileModel.FilterPattern) *filter.Profile {
	profile := &filter.Profile{
		ID:   p.ID,
		Name: p.Name,
	}

	for _, pattern := range patterns {
		switch pattern.Kind {
		case filterProfileModel.KindSearch:
			profile.Search = append(profile.Search, pattern.Pattern)
		case filterProfileModel.KindRequired:
			profile.Required = append(profile.Required, pattern.Pattern)
		case filterProfileModel.KindRejected:
			profile.Rejected = append(profile.Rejected, pattern.Pattern)
		case filterProfileModel.KindPreferred:
			profile.Preferred = append(profile.Preferred, filter.PreferredPattern{
				Pattern: pattern.Pattern,
				Score:   pattern.Score,
			})
		}
	}

	return profile
}

func BuildContext(db *sqlx.DB, show showModel.Show, season, episode int, settings EffectiveSettings) (filter.Context, error) {
	ctx := filter.Context{
		Show:    utils.NormalizeTitle(show.Name),
		Season:  season,
		Episode: episode,
	}

	if !settings.UseAliases {
		return ctx, nil
	}

	aliasesRepo := showAliasRepository.NewShowAliasesRepository(db)
	aliases, err := aliasesRepo.ListByShowID(show.ID)
	if err != nil {
		return filter.Context{}, err
	}

	seen := make(map[string]struct{})
	for _, alias := range aliases {
		if settings.OnlyLatin && !latinRegex.MatchString(alias.Alias) {
			continue
		}

		normalized := utils.NormalizeTitle(alias.Alias)
		if normalized == "" || normalized == ctx.Show {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		ctx.Aliases = append(ctx.Aliases, normalized)
	}

	return ctx, nil
}

// SearchPatterns returns the search patterns to use for a profile, falling
// back to the default pattern when the profile has none.
func SearchPatterns(profile *filter.Profile) []string {
	if profile != nil && len(profile.Search) > 0 {
		return profile.Search
	}
	return []string{filter.DefaultSearchPattern}
}
