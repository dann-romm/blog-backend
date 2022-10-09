package entity

type Tag struct {
	Id          int    `db:"id"`
	Description string `db:"description"`
}
