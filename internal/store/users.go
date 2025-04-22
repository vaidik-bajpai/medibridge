package store

import (
	"context"
	"errors"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

var (
	ErrEmailExists   = errors.New("email already exists")
	ErrUsernameTaken = errors.New("username already taken")
)

type User struct {
	client *db.PrismaClient
}

func (s *User) Create(ctx context.Context, req *dto.SignupReq) error {
	_, err := s.client.User.CreateOne(
		db.User.Username.Set(req.Username),
		db.User.Email.Set(req.Email),
		db.User.Password.Set(req.Password),
	).Exec(ctx)
	if err != nil {
		if info, ok := db.IsErrUniqueConstraint(err); ok {
			switch {
			case info.Fields[0] == db.User.Email.Field():
				return ErrEmailExists
			case info.Fields[0] == db.User.Username.Field():
				return ErrUsernameTaken
			}
		}
		return err
	}

	return nil
}
