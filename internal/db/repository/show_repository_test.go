package repository

import (
	"github.com/jusoaresg/gorgon/testutils"
	"math/rand"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func getShowRepo() (*ShowRepository, *sqlx.DB) {
	db := testutils.GetTestDB()
	showRepo := NewShowRepository(db)
	return showRepo, db
}

func TestShowRepository_Create_Success(t *testing.T) {
	showRepo, _ := getShowRepo()

	show := testutils.MakeFakeShow()

	id, err := showRepo.Create(show)
	assert.NoError(t, err)
	assert.NotZero(t, id)
}

func TestShowRepository_GetByID_Found(t *testing.T) {
	showRepo, _ := getShowRepo()

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	resultShow, err := showRepo.GetById(id)
	assert.NoError(t, err)

	show.ID = id
	assert.EqualValues(t, show, resultShow)
}

func TestShowRepository_GetByID_NotFound(t *testing.T) {
	showRepo, _ := getShowRepo()

	resultShow, err := showRepo.GetById(rand.Int63())

	assert.Empty(t, resultShow)
	assert.Error(t, err)
}

func TestShowRepository_GetByTvMazeID_Found(t *testing.T) {
	showRepo, _ := getShowRepo()

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	resultShow, err := showRepo.GetByTvMazeID(show.TvMazeID)
	assert.NoError(t, err)

	show.ID = id
	assert.EqualValues(t, show, resultShow)
}

func TestShowRepository_GetByTvMazeID_NotFound(t *testing.T) {
	showRepo, _ := getShowRepo()

	resultShow, err := showRepo.GetByTvMazeID(rand.Int63())

	assert.Empty(t, resultShow)
	assert.Error(t, err)
}

func TestShowRepository_Delete_Success(t *testing.T) {
	showRepo, _ := getShowRepo()
	show := testutils.MakeFakeShow()

	var err error
	show.ID, err = showRepo.Create(show)
	assert.NoError(t, err)

	resultShow, err := showRepo.GetById(show.ID)
	assert.NotEmpty(t, resultShow)
	assert.EqualValues(t, show, resultShow)

	err = showRepo.DeleteById(show.ID)
	assert.NoError(t, err)

	resultShowDelete, err := showRepo.GetById(show.ID)
	assert.Empty(t, resultShowDelete)
}

func TestShowRepository_Delete_NotFound(t *testing.T) {
	showRepo, _ := getShowRepo()

	err := showRepo.DeleteById(rand.Int63())
	assert.Error(t, err)
}

func TestShowRepository_Delete_Cascade(t *testing.T) {
	showRepo, db := getShowRepo()
	epRepo := NewEpisodeRepository(db)
	seasonRepo := NewSeasonRepository(db)

	show := testutils.MakeFakeShow()

	showID, err := showRepo.Create(show)
	show.ID = showID
	assert.NoError(t, err)

	season := testutils.MakeFakeSeason()
	season.ShowID = show.ID
	seasonID, err := seasonRepo.Create(season)
	season.ID = seasonID
	assert.NoError(t, err)

	episode := testutils.MakeFakeEpisode()
	episode.ShowID = show.ID
	episode.SeasonID = season.ID
	episodeID, err := epRepo.Create(episode)
	episode.ID = episodeID
	assert.NoError(t, err)

	err = showRepo.DeleteById(show.ID)
	assert.NoError(t, err)

	resultShow, _ := showRepo.GetById(show.ID)
	assert.Empty(t, resultShow)

	resultSeason, _ := seasonRepo.GetById(season.ID)
	assert.Empty(t, resultSeason)

	resultEpisode, _ := epRepo.GetByID(episode.ID)
	assert.Empty(t, resultEpisode)
}

func TestShowRepository_List_ReturnAllShows(t *testing.T) {
	showRepo, _ := getShowRepo()

	show1 := testutils.MakeFakeShow()
	show2 := testutils.MakeFakeShow()
	show3 := testutils.MakeFakeShow()
	show1.ID, _ = showRepo.Create(show1)
	show2.ID, _ = showRepo.Create(show2)
	show3.ID, _ = showRepo.Create(show3)

	shows, err := showRepo.List()
	ids := []int64{shows[0].ID, shows[1].ID, shows[2].ID}
	assert.NoError(t, err)
	assert.Contains(t, ids, show1.ID)
	assert.Contains(t, ids, show2.ID)
	assert.Contains(t, ids, show3.ID)
}

func TestShowRepository_List_ReturnEmpty(t *testing.T) {
	showRepo, _ := getShowRepo()

	shows, err := showRepo.List()
	assert.NoError(t, err)
	assert.Empty(t, shows)
}

func TestShowRepository_Update_Succes(t *testing.T) {
	showRepo, _ := getShowRepo()

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	newShow := testutils.MakeFakeShow()
	newShow.TvMazeID = show.TvMazeID
	newShow.ID = id

	err = showRepo.UpdateByTvMazeID(newShow)
	assert.NoError(t, err)

	resultShow, err := showRepo.GetByTvMazeID(newShow.TvMazeID)
	assert.EqualValues(t, newShow, resultShow)
}

func TestShowRepository_Update_NotFound(t *testing.T) {
	showRepo, _ := getShowRepo()

	show := testutils.MakeFakeShow()

	err := showRepo.UpdateByTvMazeID(show)
	assert.Error(t, err)

	resultShow, err := showRepo.GetByTvMazeID(show.TvMazeID)
	assert.NotEqualValues(t, show, resultShow)
}
