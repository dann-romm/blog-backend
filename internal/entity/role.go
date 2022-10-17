package entity

type Role struct {
	Id   int      `db:"id"`
	Name RoleType `db:"name"`
}

type RoleType string

// TODO: add support for guest role
// RoleGuest     RoleType = "guest"     // read only
const (
	RoleUser      RoleType = "user"      // can create articles and comments
	RoleModerator RoleType = "moderator" // can edit and delete articles and comments
	RoleAdmin     RoleType = "admin"     // can edit and delete users, articles and comments
)
