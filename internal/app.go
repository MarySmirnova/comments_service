package internal

import (
	"net/http"

	"github.com/MarySmirnova/comments_service/internal/api"
	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/MarySmirnova/comments_service/internal/database/postgres"

	log "github.com/sirupsen/logrus"
)

type Application struct {
	db  *postgres.Store
	cfg config.Application
}

func NewApplication(cfg config.Application) (*Application, error) {
	a := &Application{
		cfg: cfg,
	}

	err := a.initDatabase()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) initDatabase() error {
	db, err := postgres.New(a.cfg.Postgres)
	if err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *Application) StartServer() {
	s := api.NewCommentsServer(a.cfg.Server, a.db)
	srv := s.GetHTTPServer()

	log.WithFields(log.Fields{"service": s.Name, "listen": srv.Addr}).Info("start server")

	err := srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.WithField("service", s.Name).WithError(err).Error("the channel raised an error")
		return
	}
}
