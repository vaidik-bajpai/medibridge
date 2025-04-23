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

	ErrNotFound = errors.New("user not registered")
)

type User struct {
	client *db.PrismaClient
}

func (s *User) Create(ctx context.Context, req *dto.SignupReq) error {
	_, err := s.client.User.CreateOne(
		db.User.Fullname.Set(req.Fullname),
		db.User.Email.Set(req.Email),
		db.User.Activated.Set(req.Activated),
		db.User.Role.Set(db.Role(req.Role)),
		db.User.Password.Set(req.Password),
	).Exec(ctx)
	if err != nil {
		if info, ok := db.IsErrUniqueConstraint(err); ok {
			switch {
			case info.Fields[0] == db.User.Email.Field():
				return ErrEmailExists
			case info.Fields[0] == db.User.Fullname.Field():
				return ErrUsernameTaken
			}
		}
		return err
	}

	return nil
}

func (s *User) FindViaEmail(ctx context.Context, email string) (*dto.UserModel, error) {
	user, err := s.client.User.FindFirst(
		db.User.Email.Equals(email),
		db.User.Email.Mode(db.QueryModeInsensitive),
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return nil, ErrNotFound
		}
		return nil, err
	}

	pass, _ := user.Password()
	oAuthID, _ := user.OauthID()
	OAuthProvider, _ := user.OauthProvider()

	res := &dto.UserModel{
		ID:            user.ID,
		Username:      user.Fullname,
		Email:         user.Email,
		Password:      pass,
		Activated:     user.Activated,
		Role:          string(user.Role),
		OAuthID:       oAuthID,
		OAuthProvider: OAuthProvider,
	}

	return res, nil
}
