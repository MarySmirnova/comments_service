package moderator

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/MarySmirnova/comments_service/internal/database"
)

type Storage interface {
	GetUnmoderatedComments() ([]*database.Comment, error)
	UpdateModeratedComments([]*database.Comment) error
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

func (m *Moderator) Start(ctx context.Context) error {
	for {
		comms, err := m.db.GetUnmoderatedComments()
		if err != nil {
			//			log.WithError(err).Error("failed to get unmoderated comments from the database")
			continue
		}

		wg := sync.WaitGroup{}

		for _, comm := range comms {
			wg.Add(1)
			go m.Check(comm, &wg)
		}

		wg.Wait()

		if err := m.db.UpdateModeratedComments(comms); err != nil {
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

func (m *Moderator) Check(comm *database.Comment, wg *sync.WaitGroup) {
	defer wg.Done()

	comm.Banned = false

	words := strings.Fields(comm.Text)
	for _, word := range words {
		if _, ok := m.forbiddenWords[word]; ok {
			comm.Banned = true
			break
		}
	}

	comm.Moderated = true
}
