package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

type handler struct {
	validate *validator.Validate
	logger   *zap.Logger
	store    *store.Store
}

func NewHandler(v *validator.Validate, l *zap.Logger, store *store.Store) *handler {
	return &handler{
		validate: v,
		logger:   l,
		store:    store,
	}
}
