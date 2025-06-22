package repository

import (
	"math/rand"
	"testing"

	"github.com/jusoaresg/gorgon/internal/episode_content/model"
	"github.com/stretchr/testify/assert"

	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/jmoiron/sqlx"
)

func getEpContentRepo() (*sqlx.DB, EpisodeContentRepository) {
	db := testutils.GetTestDB()
	repo := NewEpisodeContentRepository(db)
	return db, repo
}

func MakeEpisodeContentWithDeps(db *sqlx.DB) model.EpisodeContent {
	fakeShow := testutils.MakeFakeShow()

	showRepo := showRepository.NewShowRepository(db)
	showID, _ := showRepo.Create(fakeShow)

	fakeSeason := testutils.MakeFakeSeason()
	fakeSeason.ShowID = showID

	seasonRepo := seasonRepository.NewSeasonRepository(db)
	seasonID, _ := seasonRepo.Create(fakeSeason)

	fakeEpisode := testutils.MakeFakeEpisode()
	fakeEpisode.ShowID = showID
	fakeEpisode.SeasonID = seasonID

	episodeRepository := episodeRepository.NewEpisodeRepository(db)
	episodeID, _ := episodeRepository.Create(fakeEpisode)

	fakeEpContent := testutils.MakeFakeEpisodeContent()
	fakeEpContent.EpisodeId = episodeID

	return fakeEpContent
}

func TestEpisodeContentRepository_Create_Success(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	epContentID, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, epContentID)

	epContent.ID = epContentID
	resultContent, err := epContentRepo.GetById(epContentID)

	assert.NoError(t, err)
	assert.NotEmpty(t, resultContent)
	assert.EqualValues(t, epContent, resultContent)
}

func TestEpisodeContentRepository_GetByID_Found(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	epContentID, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, epContentID)

	epContent.ID = epContentID
	resultContent, err := epContentRepo.GetById(epContentID)

	assert.NoError(t, err)
	assert.NotEmpty(t, resultContent)
	assert.EqualValues(t, epContent, resultContent)
}

func TestEpisodeContentRepository_GetByID_NotFound(t *testing.T) {
	_, epContentRepo := getEpContentRepo()

	resultContent, err := epContentRepo.GetById(rand.Int63n(10000))

	assert.Empty(t, resultContent)
	assert.Error(t, err)
}

func TestEpisodeContentRepository_GetByEpisodeID_Found(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	epContentID, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, epContentID)

	epContent.ID = epContentID
	resultContent, err := epContentRepo.GetByEpisodeId(epContent.EpisodeId)

	assert.NoError(t, err)
	assert.NotEmpty(t, resultContent)
	assert.EqualValues(t, epContent, resultContent)
}

func TestEpisodeContentRepository_GetByEpisodeID_NotFound(t *testing.T) {
	_, epContentRepo := getEpContentRepo()

	resultContent, err := epContentRepo.GetByEpisodeId(rand.Int63n(10000))

	assert.Empty(t, resultContent)
	assert.Error(t, err)
}

func TestEpisodeContentRepository_DeleteById_Success(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	epContentID, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, epContentID)

	err = epContentRepo.DeleteById(epContentID)
	assert.NoError(t, err)

	resultContent, err := epContentRepo.GetById(epContentID)
	assert.Error(t, err)
	assert.Empty(t, resultContent)
}

func TestEpisodeContentRepository_List_ReturnsAll(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent1 := MakeEpisodeContentWithDeps(db)
	epContent2 := MakeEpisodeContentWithDeps(db)

	_, err := epContentRepo.Create(epContent1)
	assert.NoError(t, err)

	_, err = epContentRepo.Create(epContent2)
	assert.NoError(t, err)

	list, err := epContentRepo.List()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestEpisodeContentRepository_ListByEpisodeId_ReturnAllContents(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	_, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)

	list, err := epContentRepo.ListByEpisodeId(epContent.EpisodeId)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	assert.Equal(t, epContent.EpisodeId, list[0].EpisodeId)
}

func TestEpisodeContentRepository_Update_Success(t *testing.T) {
	db, epContentRepo := getEpContentRepo()

	epContent := MakeEpisodeContentWithDeps(db)
	epContentID, err := epContentRepo.Create(epContent)
	assert.NoError(t, err)

	epContent.ID = epContentID
	epContent.Name = "Updated Name"

	err = epContentRepo.Update(epContent)
	assert.NoError(t, err)

	updatedContent, err := epContentRepo.GetById(epContentID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedContent.Name)
}
