package routes

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/episode"
	episodeApi "github.com/jusoaresg/gorgon/internal/episode/api"
	"github.com/jusoaresg/gorgon/internal/indexer"
	indexerApi "github.com/jusoaresg/gorgon/internal/indexer/api"
	"github.com/jusoaresg/gorgon/internal/season"
	seasonApi "github.com/jusoaresg/gorgon/internal/season/api"
	"github.com/jusoaresg/gorgon/internal/show"
	showApi "github.com/jusoaresg/gorgon/internal/show/api"
	"github.com/jusoaresg/gorgon/internal/show_aliases"
	showAliasesApi "github.com/jusoaresg/gorgon/internal/show_aliases/api"

	"github.com/labstack/echo/v4"
)

func SetupDatabaseRouter(v1 *echo.Group, deps *show.Dependencies, episodeDeps *episode.Dependencies, seasonDeps *season.Dependencies, indexerDeps *indexer.Dependencies, showAliasesDeps *show_aliases.Dependencies) {
	logger := config.GetLogger()

	showHandler := showApi.NewHandler(deps)
	episodeHandler := episodeApi.NewHandler(episodeDeps)
	seasonHandler := seasonApi.NewHandler(seasonDeps)
	indexerHandler := indexerApi.NewHandler(indexerDeps)
	showAliasesHandler := showAliasesApi.NewHandler(showAliasesDeps)

	listGroup := v1.Group("database/")
	{
		showGroup := listGroup.Group("show")
		{
			showGroup.POST("", showHandler.AddShowToList)
			logger.Info("POST route added to /api/v1/database/show")

			showGroup.GET("", showHandler.ListShows)
			logger.Info("GET route added to /api/v1/database/show")

			showGroup.GET("/full", showHandler.ListFullShows)
			logger.Info("GET route added to /api/v1/database/show/full")

			showGroup.GET("/:id", showHandler.GetShow)
			logger.Info("GET route added to /api/v1/database/show/:id")

			showGroup.DELETE("/:id", showHandler.DeleteShow)
			logger.Info("DELETE route added to /api/v1/database/show/:id")

			showGroup.POST("/update-info", showHandler.UpdateShowInfo)
			logger.Info("POST route added to /api/v1/database/show/update-info")

			episodeGroup := showGroup.Group("/episode")
			{
				episodeGroup.POST("/status", episodeHandler.ChangeEpisodeStatus)
				logger.Info("POST route added to /api/v1/database/show/episode/status")

				episodeGroup.POST("/search", episodeHandler.SearchProcessEpisode)
				logger.Info("POST route added to /api/v1/database/show/episode/search")

				episodeGroup.POST("/search/all", episodeHandler.SearchProcessShowWantedEpisodes)
				logger.Info("POST route added to /api/v1/database/show/episode/search/all")

				episodeGroup.POST("/search/season/:id", episodeHandler.SearchProcessSeasonEpisodes)
				logger.Info("POST route added to /api/v1/database/show/episode/search/season/:id")

				episodeGroup.GET("/:id", episodeHandler.GetShowEpisodes)
				logger.Info("GET route added to /api/v1/database/show/episode/:id")

				episodeGroup.DELETE("/:id", episodeHandler.DeleteDownloadedEpisode)
				logger.Info("DELETE route added to /api/v1/database/show/episode/:id")
			}

			aliasGroup := showGroup.Group("/aliases")
			{
				aliasGroup.GET("/:id", showAliasesHandler.GetShowAliases)
			}

			seasonsGroup := showGroup.Group("/season")
			{
				seasonsGroup.GET("/:id", seasonHandler.GetShowSeasons)
				logger.Info("GET route added to /api/v1/database/show/season/:id")
			}
		}

		indexerGroup := listGroup.Group("indexer")
		{
			indexerGroup.GET(":id", indexerHandler.GetIndexer)
			logger.Info("GET route added to /api/v1/database/indexer/:id")
			indexerGroup.GET("", indexerHandler.ListIndexers)
			logger.Info("GET route added to /api/v1/database/indexer")

			indexerGroup.POST("", indexerHandler.AddIndexer)
			logger.Info("POST route added to /api/v1/database/indexer")

			indexerGroup.DELETE("", indexerHandler.DeleteIndexer)
			logger.Info("DELETE route added to /api/v1/database/indexer")
		}
	}
	logger.Info("Database routes added successfully", slog.String("endpoint", "/api/v1/database"))
}
