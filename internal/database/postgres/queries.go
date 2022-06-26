package postgres

import (
	"errors"

	"github.com/MarySmirnova/comments_service/internal/database"
	"github.com/jackc/pgx/v4"
)

func (s *Store) NewComment(comm *database.Comment) error {
	query := `
	INSERT INTO comments.posts (
		parent_id,
		news_id,
		text,
		pub_time)
	VALUES ($1, $2, $3, $4);`

	_, err := s.db.Exec(ctx, query, comm.ParentID, comm.NewsID, comm.Text, comm.PubTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAllCommentsByNewsID(newsID int) ([]*database.Comment, error) {
	query := `
	SELECT 
		id,
		parent_id,
		news_id,
		text,
		pub_time
	FROM comments.posts
	WHERE news_id = $1
	ORDER BY pub_time DESC;`

	var comments []*database.Comment

	rows, err := s.db.Query(ctx, query, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comm database.Comment

		err = rows.Scan(&comm.ID, &comm.ParentID, &comm.NewsID, &comm.Text, &comm.PubTime)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comm)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return comments, nil
}
