package service

import (
	"github.com/jusoaresg/gorgon/config"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showAliasModel "github.com/jusoaresg/gorgon/internal/show_aliases/model"
	"github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/utils"
)

func GetTitleAlias(show showModel.Show) ([]showAliasModel.ShowAlias, error) {
	db := config.GetSQLite()
	aliasRepo := repository.NewShowAliasesRepository(db)

	aliases, err := aliasRepo.ListByShowID(show.ID)
	if err != nil {
		return nil, err
	}
	return aliases, err
}

func GetNormalizedTitleAlias(show showModel.Show) ([]string, error) {
	aliases, err := GetTitleAlias(show)
	if err != nil {
		return nil, err
	}

	titleAliases := []string{
		utils.NormalizeTitle(show.Name),
	}
	for _, alias := range aliases {
		if alias.Alias == show.Name {
			continue
		}
		normalized := utils.NormalizeTitle(alias.Alias)
		titleAliases = append(titleAliases, normalized)
	}
	return titleAliases, nil
}
