package api

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/MarySmirnova/comments_service/internal/database"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextReqIDKey ContextKey = "request_id"

const itemsPerPage = 15

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
	handler.Use(s.reqIDMiddleware, s.logMiddleware)
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
	w.Header().Add("Code", strconv.Itoa(code))
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (s *Server) internalError(w http.ResponseWriter, err error) {
	log.WithField("service", s.Name).WithError(err).Error("something went wrong")
	w.Header().Add("Code", strconv.Itoa(http.StatusInternalServerError))
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("something went wrong"))
}

func (s *Server) reqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID int
		reqIDString := r.FormValue("request_id")

		if reqIDString == "" {
			reqID = s.generateReqID()
		}

		if reqIDString != "" {
			id, err := strconv.Atoi(reqIDString)
			if err != nil {
				s.writeResponseError(w, err, http.StatusBadRequest)
				return
			}
			reqID = id
		}

		ctx := context.WithValue(r.Context(), ContextReqIDKey, reqID)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.WithFields(log.Fields{
				"request_time": time.Now().Format("2006-01-02 15:04:05.000000"),
				"request_ip":   strings.TrimPrefix(strings.Split(r.RemoteAddr, ":")[1], "["),
				"code":         w.Header().Get("Code"),
				"request_id":   r.Context().Value(ContextReqIDKey),
			}).Info("news reader response")
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) generateReqID() int {
	max := 999999999999
	min := 100000

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
