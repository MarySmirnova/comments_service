package moderator

import (
	"context"
	"strings"
	"time"

	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/MarySmirnova/comments_service/internal/database"
)

type Storage interface {
	GetUnmoderatedComments() ([]*database.Comment, error)
	UpdateModeratedComments([]ModerationStatus) error
}

type Moderator struct {
	db             Storage
	checkInterval  time.Duration
	forbiddenWords map[string]struct{}
}

func New(cfg config.Moderator, db Storage) *Moderator {
	forbiddenWords := map[string]struct{}{
		"qwerty": struct{}{},
		"йцукен": struct{}{},
		"zxvbnm": struct{}{},
	}

	return &Moderator{
		db:             db,
		checkInterval:  cfg.CheckInterval,
		forbiddenWords: forbiddenWords,
	}
}

type ModerationStatus struct {
	comment *database.Comment
	blocked bool
}

func (m *Moderator) Start(ctx context.Context) error {
	var chanStatus = make(chan ModerationStatus)

	for {
		var moderated []ModerationStatus

		comms, err := m.db.GetUnmoderatedComments()
		if err != nil {
			//			log.WithError(err).Error("failed to get unmoderated comments from the database")
			continue
		}

		for _, comm := range comms {
			go m.Check(comm, chanStatus)
		}

		checkCount := len(comms)
		for checkCount > 0 {
			select {
			case comm := <-chanStatus:
				moderated = append(moderated, comm)
			}
		}

		if err := m.db.UpdateModeratedComments(moderated); err != nil {
			//			log.WithError(err).Error("failed to update moderated comments in the database")
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.checkInterval):
		}
	}
}

func (m *Moderator) Check(comm *database.Comment, chanStatus chan<- ModerationStatus) {
	status := ModerationStatus{
		comment: comm,
		blocked: false,
	}

	words := strings.Fields(comm.Text)
	for _, word := range words {
		if _, ok := m.forbiddenWords[word]; ok {
			status.blocked = true
			chanStatus <- status
			return
		}
	}

	chanStatus <- status
}
