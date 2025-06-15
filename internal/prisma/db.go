package database

import (
	"log"

	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

func NewPrismaClient() (*db.PrismaClient, error) {
	client := db.NewClient()

	if err := client.Connect(); err != nil {

		log.Println(err)
		return nil, err
	}

	return client, nil
}
