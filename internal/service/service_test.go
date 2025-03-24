package service_test

import (
	"context"
	"errors"
	"testing"
	"users-app/internal/entity"
	"users-app/internal/mocks"
	"users-app/internal/service"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_GetUserByID(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.New(mockRepo)

	ctx := context.Background()

	userID := uuid.Must(uuid.NewV4())
	repositoryErr := errors.New("repository error")

	tests := []struct {
		name         string
		userID       uuid.UUID
		expectedUser entity.User
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:         "User found",
			userID:       userID,
			expectedUser: entity.User{ID: userID, Name: "test"},
			expectedErr:  nil,
			mockBehavior: func() {
				mockRepo.EXPECT().GetUserByID(ctx, gomock.Any()).Return(
					entity.User{
						ID:   userID,
						Name: "test"}, nil,
				)
			},
		},
		{
			name:         "User not found",
			userID:       uuid.Must(uuid.NewV4()),
			expectedUser: entity.User{},
			expectedErr:  entity.ErrNotFound,
			mockBehavior: func() {
				mockRepo.EXPECT().GetUserByID(ctx, gomock.Any()).Return(entity.User{}, entity.ErrNotFound)
			},
		},
		{
			name:         "Repository error",
			userID:       uuid.Must(uuid.NewV4()),
			expectedUser: entity.User{},
			expectedErr:  repositoryErr,
			mockBehavior: func() {
				mockRepo.EXPECT().GetUserByID(ctx, gomock.Any()).Return(entity.User{}, repositoryErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			user, err := svc.GetUserByID(ctx, tt.userID)
			if tt.expectedErr != nil {
				r.Error(err)
				r.ErrorIs(err, tt.expectedErr)
			} else {
				r.NoError(err)
				r.Equal(tt.expectedUser, user)
			}
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.New(mockRepo)

	ctx := context.Background()

	user := entity.User{ID: uuid.Must(uuid.NewV4()), Name: "test"}
	repositoryErr := errors.New("repository error")

	tests := []struct {
		name         string
		user         entity.User
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:        "Create user successfully",
			user:        user,
			expectedErr: nil,
			mockBehavior: func() {
				mockRepo.EXPECT().CreateUser(ctx, user).Return(nil)
			},
		},
		{
			name:        "User with email already exists",
			user:        user,
			expectedErr: entity.ErrAlreadyExists,
			mockBehavior: func() {
				mockRepo.EXPECT().CreateUser(ctx, user).Return(entity.ErrAlreadyExists)
			},
		},
		{
			name:        "Repository error",
			user:        user,
			expectedErr: repositoryErr,
			mockBehavior: func() {
				mockRepo.EXPECT().CreateUser(ctx, user).Return(repositoryErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := svc.CreateUser(ctx, tt.user)
			if tt.expectedErr != nil {
				r.Error(err)
				r.ErrorIs(err, tt.expectedErr)
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.New(mockRepo)

	ctx := context.Background()

	user := entity.User{ID: uuid.Must(uuid.NewV4()), Name: "test"}
	repositoryErr := errors.New("repository error")

	tests := []struct {
		name         string
		user         entity.User
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:        "Update user successfully",
			user:        user,
			expectedErr: nil,
			mockBehavior: func() {
				mockRepo.EXPECT().UpdateUser(ctx, user).Return(nil)
			},
		},
		{
			name:        "User not found",
			user:        user,
			expectedErr: entity.ErrNotFound,
			mockBehavior: func() {
				mockRepo.EXPECT().UpdateUser(ctx, user).Return(entity.ErrNotFound)
			},
		},
		{
			name:        "Repository error",
			user:        user,
			expectedErr: repositoryErr,
			mockBehavior: func() {
				mockRepo.EXPECT().UpdateUser(ctx, user).Return(repositoryErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := svc.UpdateUser(ctx, tt.user)
			if tt.expectedErr != nil {
				r.Error(err)
				r.ErrorIs(err, tt.expectedErr)
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.New(mockRepo)

	ctx := context.Background()

	userID := uuid.Must(uuid.NewV4())
	repositoryErr := errors.New("repository error")

	tests := []struct {
		name         string
		userID       uuid.UUID
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:        "Delete user successfully",
			userID:      userID,
			expectedErr: nil,
			mockBehavior: func() {
				mockRepo.EXPECT().DeleteUser(ctx, userID).Return(nil)
			},
		},
		{
			name:        "User not found",
			userID:      userID,
			expectedErr: entity.ErrNotFound,
			mockBehavior: func() {
				mockRepo.EXPECT().DeleteUser(ctx, userID).Return(entity.ErrNotFound)
			},
		},
		{
			name:        "Repository error",
			userID:      userID,
			expectedErr: repositoryErr,
			mockBehavior: func() {
				mockRepo.EXPECT().DeleteUser(ctx, userID).Return(repositoryErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := svc.DeleteUser(ctx, tt.userID)
			if tt.expectedErr != nil {
				r.Error(err)
				r.ErrorIs(err, tt.expectedErr)
			} else {
				r.NoError(err)
			}
		})
	}
}




