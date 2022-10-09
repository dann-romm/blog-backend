package entity

import "time"

type Article struct {
	Id             int       `db:"id"`
	AuthorId       int       `db:"author_id"`
	Title          string    `db:"title"`
	Description    string    `db:"description"`
	Content        string    `db:"content"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	CommentsCount  int       `db:"comments_count"`
	FavoritesCount int       `db:"favorites_count"`
	ViewsCount     int       `db:"views_count"`
	VotesUpCount   int       `db:"votes_up_count"`
	VotesDownCount int       `db:"votes_down_count"`
}
