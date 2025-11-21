package repository

import (
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserInterface interface {
	Create(u *model.User) (bool, error)
	GetByEmail(email string) (*model.User, error)
	UpdatePassword(email string, password string) (bool, error)
	UpdateUsername(email string, name string) (bool, error)
	DeleteUser(email string) (bool, error)
	ListAll() ([]model.User, error)
}
