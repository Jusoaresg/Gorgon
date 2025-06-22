package repository

import (
	"github.com/jusoaresg/gorgon/internal/db/model"
	"github.com/jusoaresg/gorgon/testutils"
	"math/rand"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func MakeSeasonWithDeps(db *sqlx.DB) (model.Season, error) {
	showRepo := NewShowRepository(db)

	show := testutils.MakeFakeShow()
	showID, _ := showRepo.Create(show)

	season := testutils.MakeFakeSeason()
	season.ShowID = showID

	return season, nil
}

func getSeasonRepo() (*SeasonRepository, *sqlx.DB) {
	db := testutils.GetTestDB()
	seasonRepo := NewSeasonRepository(db)
	return seasonRepo, db
}

func TestSeasonRepository_Create_Success(t *testing.T) {
	seasonRepo, db := getSeasonRepo()

	season, err := MakeSeasonWithDeps(db)

	id, err := seasonRepo.Create(season)
	assert.NoError(t, err)
	assert.NotZero(t, id)
}

func TestSeasonRepository_GetByID_Found(t *testing.T) {
	seasonRepo, db := getSeasonRepo()
	var err error

	season, _ := MakeSeasonWithDeps(db)
	season.ID, err = seasonRepo.Create(season)
	assert.NoError(t, err)

	resultEp, _ := seasonRepo.GetById(season.ID)
	assert.NoError(t, err)

	assert.EqualValues(t, season, resultEp)
}

func TestSeasonRepository_GetByID_NotFound(t *testing.T) {
	seasonRepo, _ := getSeasonRepo()

	resultEp, err := seasonRepo.GetById(rand.Int63())

	assert.Empty(t, resultEp)
	assert.Error(t, err)
}

func TestSeasonRepository_List_ReturnAllSeasons(t *testing.T) {
	seasonRepo, db := getSeasonRepo()

	season1, _ := MakeSeasonWithDeps(db)
	season2, _ := MakeSeasonWithDeps(db)
	season3, _ := MakeSeasonWithDeps(db)

	season1.ID, _ = seasonRepo.Create(season1)
	season2.ID, _ = seasonRepo.Create(season2)
	season3.ID, _ = seasonRepo.Create(season3)

	seasons, err := seasonRepo.List()
	ids := []int64{seasons[0].ID, seasons[1].ID, seasons[2].ID}
	assert.NoError(t, err)
	assert.Contains(t, ids, season1.ID)
	assert.Contains(t, ids, season2.ID)
	assert.Contains(t, ids, season3.ID)
}

func TestSeasonRepository_List_ReturnEmpty(t *testing.T) {
	seasonRepo, _ := getSeasonRepo()

	seasons, err := seasonRepo.List()
	assert.NoError(t, err)
	assert.Empty(t, seasons)
}
