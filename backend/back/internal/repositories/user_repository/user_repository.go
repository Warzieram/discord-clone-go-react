package userrepository

import (
	"back/internal/models/user"
	"context"
)

type UserRepository interface {
	Save(ctx *context.Context, u *user.User) error
	GetByID(ctx *context.Context, id int) (user.User, error)
	GetByEmail(ctx *context.Context, email string) (user.User, error)
}
