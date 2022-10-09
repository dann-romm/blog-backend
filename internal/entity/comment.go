package entity

import "time"

type Comment struct {
	Id             int       `db:"id"`
	AuthorId       int       `db:"author_id"`
	ArticleId      int       `db:"article_id"`
	ParentId       *int      `db:"parent_id"`
	Content        string    `db:"content"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	VotesUpCount   int       `db:"votes_up_count"`
	VotesDownCount int       `db:"votes_down_count"`
}
