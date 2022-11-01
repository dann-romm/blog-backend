package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                     uuid.UUID `db:"id"`
	Name                   string    `db:"name"`
	Username               string    `db:"username"`
	Password               string    `db:"password"`
	Email                  string    `db:"email"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
	Role                   RoleType  `db:"role"`
	Description            string    `db:"description"`
	ArticlesCount          int       `db:"articles_count"`
	CommentsCount          int       `db:"comments_count"`
	FavoritesArticlesCount int       `db:"favorites_articles_count"`
	FavoritesCommentsCount int       `db:"favorites_comments_count"`
	FollowersCount         int       `db:"followers_count"`
	FollowingCount         int       `db:"following_count"`
}

type RoleType string

const (
	RoleGuest     RoleType = "guest"     // read only
	RoleUser      RoleType = "user"      // can create,  articles and comments, can vote
	RoleModerator RoleType = "moderator" // can delete articles and comments
	RoleAdmin     RoleType = "admin"     // can edit and delete users, articles and comments
)
