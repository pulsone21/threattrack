package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	CreatedAt uint      `json:"created_at"`
	//TODO Role Concecpt
	//TODO Auth Concept
}

func NewUser(firstname, lastname, email string) *User {
	epoch := time.Now().Unix()

	return &User{
		Id:        uuid.New(),
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Fullname:  fmt.Sprintf("%s %s", firstname, lastname),
		CreatedAt: uint(epoch),
	}
}

func (u *User) ScanTo(scan ScanFunc) error {
	return scan(
		&u.Id,
		&u.Firstname,
		&u.Lastname,
		&u.Email,
		&u.Fullname,
		&u.CreatedAt)
}
