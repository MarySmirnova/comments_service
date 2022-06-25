package database

import "errors"

var ErrNewsIDNotExist error = errors.New("no comments found for this post")

type Comment struct {
	ID        int    // id комментария
	ParentID  int    // id родительского комментария
	NewsID    int    // id новости
	Text      string // тело комментария
	PubTime   int64  // время публикации
	Moderated bool   // прошел ли модерацию
	Banned    bool   // заблокирован ли модератором
}
