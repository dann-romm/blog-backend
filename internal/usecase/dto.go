package usecase

import "github.com/google/uuid"

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type AuthParseTokenInput struct {
	Token string
}

type UserCreateUserInput struct {
	Name     string
	Username string
	Password string
	Email    string
}

type UserUpdateUserPasswordInput struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

type UserUpdateUserEmailInput struct {
	UserID uuid.UUID
	Email  string
}

type ArticleCreateArticleInput struct {
	AuthorID    uuid.UUID
	Title       string
	Description string
	Content     string
}

type ArticleGetArticleByIDInput struct {
	ID uuid.UUID
}

type ArticleGetArticlesByAuthorIDInput struct {
	AuthorID uuid.UUID
}

type ArticleGetNewestArticlesInput struct {
	Limit  int
	Offset int
}

type ArticleSetArticleFavoriteInput struct {
	UserID    uuid.UUID
	ArticleID uuid.UUID
}

type ArticleRemoveArticleFavoriteInput struct {
	UserID    uuid.UUID
	ArticleID uuid.UUID
}

type ArticleGetFavoriteArticlesInput struct {
	UserID uuid.UUID
}
