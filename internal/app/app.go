package app

import (
	"avito_bootcamp/configs"
	"avito_bootcamp/internal/controller/http"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/internal/usecase/repo"
	"avito_bootcamp/pkg/httpserver"
	"avito_bootcamp/pkg/logger"
	"avito_bootcamp/pkg/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *configs.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	flatUseCase := usecase.NewFlatUseCase(
		repo.NewFlatRepo(pg))
	houseUseCase := usecase.NewHouseUseCase(
		repo.NewHouseRepo(pg))
	authUseCase := usecase.NewAuthUseCase(repo.NewAuthRepo(pg), cfg.Token)

	handler := gin.New()
	http.NewRouter(handler, l, cfg.Token, flatUseCase, houseUseCase, authUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
