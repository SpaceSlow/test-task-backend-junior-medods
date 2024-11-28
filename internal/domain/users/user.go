package users

import "github.com/google/uuid"

type User struct {
	id               uuid.UUID
	email            string
	refreshTokenHash string
}

func NewUser(id uuid.UUID, email string, refreshTokenHash string) *User {
	return &User{
		id:               id,
		email:            email,
		refreshTokenHash: refreshTokenHash,
	}
}

func (u User) Id() uuid.UUID {
	return u.id
}

func (u User) RefreshTokenHash() string {
	return u.refreshTokenHash
}
