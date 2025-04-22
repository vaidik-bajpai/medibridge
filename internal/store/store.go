package store

import (
	"context"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type UserStorer interface {
	Create(context.Context, *dto.SignupReq) error
}

type Store struct {
	User UserStorer
}

func NewStore(client *db.PrismaClient) *Store {
	return &Store{
		User: &User{client: client},
	}
}
