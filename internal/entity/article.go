package entity

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id             uuid.UUID `db:"id"`
	AuthorId       uuid.UUID `db:"author_id"`
	Title          string    `db:"title"`
	Description    string    `db:"description"`
	Content        string    `db:"content"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	ViewsCount     int       `db:"views_count"`
	CommentsCount  int       `db:"comments_count"`
	FavoritesCount int       `db:"favorites_count"`
	VotesUpCount   int       `db:"votes_up_count"`
	VotesDownCount int       `db:"votes_down_count"`
}
