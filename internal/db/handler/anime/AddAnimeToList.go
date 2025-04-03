package anime

import (
	"gorgon/config"
	"gorgon/external/anilist/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/schemas/dtos"
	"gorgon/pkg/services"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Add Anime
// @Description Add anime to list
// @Tags Database/Anime
// @Accept json
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [post]
func AddAnimeToList(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to Add Anime To List", slog.String("endpoint", "/api/v1/database/anime"), slog.String("method", "POST"))

	var request schemas.IdRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind body request")
		schemas.SendError(c, 500, "Failed to bind body request")
		return
	}

	anilistService := service.NewAnimeListService(logger)
	response, err := anilistService.GetAnimeInfoById(request.Id)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return
	}

	anime := GroupSeasons(*response, logger)
	if err != nil {
		return
	}

	baseService := services.NewBaseService()
	if err := baseService.Add(&anime); err != nil {
		logger.Error("Failed to add anime to database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add anime to database")
		return
	}

	schemas.SendSucess(c, "Add Anime To List", &anime)
}

func GroupSeasons(anime dtos.AnimeDto, logger *slog.Logger) model.Anime {
	logger.Info("Starting group seasons process", slog.Any("anime_dto", anime))

	mainAnime := model.Anime{
		Aid:          anime.Id,
		EpisodeCount: anime.Episodes,
		Description:  anime.Description,
		Genres:       strings.Join(anime.Genres, ","),
		BannerImage:  anime.BannerImage,
		CoverImage:   anime.CoverImage.ExtraLarge,
		Status:       anime.Status,
		Title: model.Title{
			Aid:     anime.Id,
			English: anime.Title.English,
			Romaji:  anime.Title.Romaji,
		},
		Seasons:   []model.Season{},
		Relations: model.RelationWrapper{},
	}

	if anime.Relations.Edges == nil || len(anime.Relations.Edges) == 0 {
		logger.Warn("No relations found for anime", slog.Int("anime_id", anime.Id))
		return mainAnime
	}

	seasonNumber := 1 // Começamos com a primeira temporada

	for _, edge := range anime.Relations.Edges {

		logger.Info("Processing edge",
			slog.Int("relation_id", edge.Id),
			slog.Int("anime_id", edge.Id),
			slog.String("relation_type", edge.RelationType),
			slog.String("title_romaji", edge.Node.Title.Romaji),
			slog.String("title_english", edge.Node.Title.English),
			slog.String("format", edge.Node.Format),
			slog.Int("episodes", edge.Node.Episodes),
			slog.String("status", edge.Node.Status),
		)

		switch edge.RelationType {
		case "SEQUEL", "PREQUEL":

			logger.Info("Adding as season",
				slog.Int("season_number", seasonNumber),
				slog.String("title_romaji", edge.Node.Title.Romaji),
			)

			season := model.Season{
				Aid:          edge.Id,
				RomajiTitle:  edge.Node.Title.Romaji,
				EnglishTitle: edge.Node.Title.English,
				SeasonNumber: seasonNumber,
			}
			mainAnime.Seasons = append(mainAnime.Seasons, season)
			seasonNumber++

		case "SIDE_STORY", "MOVIE", "OVA":
			logger.Info("Adding as related content",
				slog.String("type", edge.RelationType),
				slog.String("title_romaji", edge.Node.Title.Romaji),
			)
			special := model.Related{
				Aid: edge.Id,
				Title: model.Title{
					Aid:     edge.Id,
					English: edge.Node.Title.English,
					Romaji:  edge.Node.Title.Romaji,
				},
				Episodes: edge.Node.Episodes,
				Type:     edge.RelationType,
				Format:   edge.Node.Format,
				Status:   edge.Node.Status,
			}
			mainAnime.Relations.Edges = append(mainAnime.Relations.Edges, special)
		}
	}

	if len(mainAnime.Seasons) == 0 {
		logger.Warn("No SEQUEL or PREQUEL found, setting default season", slog.Int("anime_id", anime.Id))
		mainAnime.Seasons = append(mainAnime.Seasons, model.Season{
			Aid:          anime.Id,
			SeasonNumber: 1,
		})
	}

	return mainAnime
}
