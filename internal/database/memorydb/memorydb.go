package memorydb

import (
	"sync"

	"github.com/MarySmirnova/comments_service/internal/database"
)

type MemoryDB struct {
	mu       sync.RWMutex
	comments map[int][]*database.Comment
	lastID   int
}

func New() *MemoryDB {
	return &MemoryDB{
		mu:       sync.RWMutex{},
		comments: make(map[int][]*database.Comment),
		lastID:   0,
	}
}

func (db *MemoryDB) NewComment(comm *database.Comment) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.lastID++
	comm.ID = db.lastID

	db.comments[comm.NewsID] = append(db.comments[comm.NewsID], comm)

	return nil
}

func (db *MemoryDB) GetAllCommentsByNewsID(newsID int) ([]*database.Comment, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	comms, ok := db.comments[newsID]
	if !ok {
		return nil, database.ErrNewsIDNotExist
	}

	return comms, nil
}
