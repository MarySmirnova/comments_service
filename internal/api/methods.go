package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MarySmirnova/comments_service/internal/database"
	"github.com/gorilla/mux"
)

func (s *Server) NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.writeResponseError(w, fmt.Errorf("invalid parameter passed: %s", err), http.StatusBadRequest)
		return
	}

	parentID, err := strconv.Atoi(r.FormValue("comm_id"))
	if err != nil {
		parentID = 0
	}

	var comment *database.Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		s.writeResponseError(w, fmt.Errorf("wrong JSON: %s", err), http.StatusBadRequest)
		return
	}

	comment.NewsID = newsID
	comment.ParentID = parentID
	comment.PubTime = time.Now().Unix()

	if err := s.db.NewComment(comment); err != nil {
		s.internalError(w, err)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusNoContent))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) AllCommentsHandler(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.writeResponseError(w, fmt.Errorf("invalid parameter passed: %s", err), http.StatusBadRequest)
		return
	}

	comments, err := s.db.GetAllCommentsByNewsID(newsID)
	if err != nil {
		if errors.Is(err, database.ErrNewsIDNotExist) {
			s.writeResponseError(w, fmt.Errorf("news id %d does not exist: %s", newsID, err), http.StatusBadRequest)
			return
		}
		s.internalError(w, err)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}
