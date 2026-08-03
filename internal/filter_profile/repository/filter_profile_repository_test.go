package repository

import (
	"testing"

	"github.com/jusoaresg/gorgon/internal/filter_profile/model"
	"github.com/jusoaresg/gorgon/testutils"
	"github.com/stretchr/testify/assert"
)

func getFilterProfileRepo() FilterProfileRepository {
	return NewFilterProfileRepository(testutils.GetTestDB())
}

func TestFilterProfileRepository_CreateAndGet(t *testing.T) {
	repo := getFilterProfileRepo()

	profile := model.FilterProfile{Name: "HD"}

	id, err := repo.Create(profile, []model.FilterPattern{
		{Kind: model.KindRequired, Pattern: "multisub"},
		{Kind: model.KindRejected, Pattern: "hdtv"},
		{Kind: model.KindPreferred, Pattern: "1080p", Score: 20},
	})
	assert.NoError(t, err)
	assert.NotZero(t, id)

	got, patterns, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "HD", got.Name)
	assert.Len(t, patterns, 3)
	assert.Equal(t, model.KindRequired, patterns[0].Kind)
	assert.Equal(t, "multisub", patterns[0].Pattern)
	assert.Equal(t, 20, patterns[2].Score)
	assert.Equal(t, 2, patterns[2].Position)
}

func TestFilterProfileRepository_GetByIDNotFound(t *testing.T) {
	repo := getFilterProfileRepo()

	_, _, err := repo.GetByID(9999)
	assert.Error(t, err)
}

func TestFilterProfileRepository_Update(t *testing.T) {
	repo := getFilterProfileRepo()

	id, err := repo.Create(model.FilterProfile{Name: "Old"}, []model.FilterPattern{
		{Kind: model.KindRequired, Pattern: "old"},
	})
	assert.NoError(t, err)

	err = repo.Update(model.FilterProfile{ID: id, Name: "New"}, []model.FilterPattern{
		{Kind: model.KindRejected, Pattern: "new"},
	})
	assert.NoError(t, err)

	got, patterns, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, "New", got.Name)
	assert.Len(t, patterns, 1)
	assert.Equal(t, model.KindRejected, patterns[0].Kind)
	assert.Equal(t, "new", patterns[0].Pattern)
}

func TestFilterProfileRepository_UpdateNotFound(t *testing.T) {
	repo := getFilterProfileRepo()

	err := repo.Update(model.FilterProfile{ID: 9999, Name: "X"}, nil)
	assert.ErrorIs(t, err, ErrFilterProfileNotFound)
}

func TestFilterProfileRepository_Delete(t *testing.T) {
	repo := getFilterProfileRepo()

	id, err := repo.Create(model.FilterProfile{Name: "ToDelete"}, []model.FilterPattern{
		{Kind: model.KindRequired, Pattern: "x"},
	})
	assert.NoError(t, err)

	err = repo.Delete(id)
	assert.NoError(t, err)

	_, _, err = repo.GetByID(id)
	assert.Error(t, err)
}

func TestFilterProfileRepository_DeleteNotFound(t *testing.T) {
	repo := getFilterProfileRepo()

	err := repo.Delete(9999)
	assert.ErrorIs(t, err, ErrFilterProfileNotFound)
}

func TestFilterProfileRepository_DeleteCascadesPatterns(t *testing.T) {
	repo := getFilterProfileRepo()

	id, err := repo.Create(model.FilterProfile{Name: "Cascade"}, []model.FilterPattern{
		{Kind: model.KindRequired, Pattern: "x"},
	})
	assert.NoError(t, err)

	err = repo.Delete(id)
	assert.NoError(t, err)

	var count int
	err = repo.db.Get(&count, "SELECT COUNT(*) FROM filter_patterns WHERE profile_id = ?", id)
	assert.NoError(t, err)
	assert.Zero(t, count)
}

func TestFilterProfileRepository_List(t *testing.T) {
	repo := getFilterProfileRepo()

	id1, err := repo.Create(model.FilterProfile{Name: "Bravo"}, nil)
	assert.NoError(t, err)
	id2, err := repo.Create(model.FilterProfile{Name: "Alpha"}, nil)
	assert.NoError(t, err)

	profiles, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)
	assert.Equal(t, "Alpha", profiles[0].Name)
	assert.Equal(t, "Bravo", profiles[1].Name)

	ids := []int64{profiles[0].ID, profiles[1].ID}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
}
