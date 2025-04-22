package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.SignupReq
		mockErr       error
		expectedError error
	}{
		{
			name: "Success",
			req: &dto.SignupReq{
				Username: "vaidik",
				Email:    "vaidikbajpai2@gmail.com",
				Password: "password123",
			},
			mockErr:       nil,
			expectedError: nil,
		},
		{
			name: "EmailAlreadyExists",
			req: &dto.SignupReq{
				Username: "vaidik",
				Email:    "vaidikbajpai2@gmail.com",
				Password: "password123",
			},
			mockErr: &types.ErrUniqueConstraint{
				Fields: []string{db.User.Email.Field()},
			}, // Simulate unique constraint error for email
			expectedError: ErrEmailExists, // Expect custom error for email already taken
		},
		{
			name: "UsernameAlreadyExists",
			req: &dto.SignupReq{
				Username: "vaidik",
				Email:    "newemail@example.com",
				Password: "password123",
			},
			mockErr:       db.ErrUniqueConstraint, // Simulate unique constraint error for username
			expectedError: ErrUsernameTaken,       // Expect custom error for username already taken
		},
		{
			name: "GeneralError",
			req: &dto.SignupReq{
				Username: "vaidik",
				Email:    "vaidikbajpai2@gmail.com",
				Password: "password123",
			},
			mockErr:       fmt.Errorf("general error"), // Simulate a general error
			expectedError: fmt.Errorf("general error"), // Expect the general error to be returned
		},
	}

	client, mock, ensure := db.NewMock()
	defer ensure(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior for CreateOne method
			if tt.mockErr != nil {
				mock.User.Expect(
					client.User.CreateOne(
						db.User.Username.Set(tt.req.Username),
						db.User.Email.Set(tt.req.Email),
						db.User.Password.Set(tt.req.Password),
					),
				).Errors(tt.mockErr)
			} else {
				expected := db.UserModel{
					InnerUser: db.InnerUser{
						Username: tt.req.Username,
						Email:    tt.req.Email,
						Password: tt.req.Password,
					},
				}
				mock.User.Expect(
					client.User.CreateOne(
						db.User.Username.Set(tt.req.Username),
						db.User.Email.Set(tt.req.Email),
						db.User.Password.Set(tt.req.Password),
					),
				).Returns(expected)
			}

			userStore := &User{client: client}

			// Run the function
			err := userStore.Create(context.Background(), tt.req)

			// Check if the error matches the expected one
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
