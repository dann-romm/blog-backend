package entity

import "github.com/google/uuid"

type Tag struct {
	Id          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}
