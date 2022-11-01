package entity

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Id             uuid.UUID     `db:"id"`
	AuthorId       uuid.UUID     `db:"author_id"`
	ArticleId      uuid.UUID     `db:"article_id"`
	ParentId       uuid.NullUUID `db:"parent_id"`
	Content        string        `db:"content"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
	VotesUpCount   int           `db:"votes_up_count"`
	VotesDownCount int           `db:"votes_down_count"`
}
