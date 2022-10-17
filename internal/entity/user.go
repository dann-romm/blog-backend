package entity

import "time"

type User struct {
	Id                     int       `db:"id"`
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
