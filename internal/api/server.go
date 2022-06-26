package api

import (
	"net/http"

	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/MarySmirnova/comments_service/internal/database"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Storage interface {
	NewComment(*database.Comment) error
	GetAllCommentsByNewsID(int) ([]*database.Comment, error)
}

type Server struct {
	Name       string
	db         Storage
	httpServer *http.Server
}

func NewCommentsServer(cfg config.Server, db Storage) *Server {
	s := &Server{
		Name: "comments",
		db:   db,
	}

	handler := mux.NewRouter()
	handler.Name("new_comment").Methods(http.MethodPost).Path("/comment/{id}").HandlerFunc(s.NewCommentHandler)
	handler.Name("get_all_comments").Methods(http.MethodGet).Path("/comment/{id}").HandlerFunc(s.AllCommentsHandler)

	s.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      handler,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	return s
}

func (s *Server) GetHTTPServer() *http.Server {
	return s.httpServer
}

func (s *Server) writeResponseError(w http.ResponseWriter, err error, code int) {
	log.WithField("service", s.Name).WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (s *Server) internalError(w http.ResponseWriter, err error) {
	log.WithField("service", s.Name).WithError(err).Error("something went wrong")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("something went wrong"))
}
