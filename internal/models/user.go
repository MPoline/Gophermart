package models

import (
	"github.com/google/uuid"
)

var users = make(map[string]*User)

type User struct {
	ID             string
	Login          string
	HashedPassword []byte
}

func NewUser(login string, hashedPassword []byte) *User {
	uniqueId, _ := uuid.NewRandom()

	newUser := &User{
		ID:             uniqueId.String(),
		Login:          login,
		HashedPassword: hashedPassword,
	}
	return newUser
}
