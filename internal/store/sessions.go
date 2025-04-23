package store

import (
	"context"
	"time"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Session struct {
	client *db.PrismaClient
}

func (s *Session) Create(ctx context.Context, req *dto.CreateSessReq) error {
	_, err := s.client.Session.CreateOne(
		db.Session.User.Link(
			db.User.ID.Equals(req.UserID),
		),
		db.Session.Token.Set(req.Token),
		db.Session.ExpiresAt.Set(req.Expiry),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) FindUserByToken(ctx context.Context, token string) (*dto.UserModel, error) {
	session, err := s.client.Session.FindFirst(
		db.Session.Token.Equals(token),
		db.Session.ExpiresAt.Gt(time.Now()),
	).With(
		db.Session.User.Fetch(),
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user := session.User()

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
