package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"users-app/internal/entity"
	"users-app/internal/mocks"
	"users-app/pkg/logger"

	"github.com/gofrs/uuid/v5"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
)

func TestHandler_GetUserByID(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)

	log, err := logger.New("mock")
	r.NoError(err)

	handler := New(log, mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockBehavior   func(userID uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().GetUserByID(gomock.Any(), userID).Return(entity.User{ID: userID, Name: "test"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing id",
			userID:         "",
			mockBehavior:   func(userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid id",
			userID:         "invalid-id",
			mockBehavior:   func(userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().GetUserByID(gomock.Any(), userID).Return(entity.User{}, entity.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "internal server error",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().GetUserByID(gomock.Any(), userID).Return(entity.User{}, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success" || tt.name == "internal server error" || tt.name == "user not found" {
				userID, err := uuid.FromString(tt.userID)
				r.NoError(err)
				tt.mockBehavior(userID)
			}

			req, err := http.NewRequest(http.MethodGet, "/user?id="+tt.userID, nil)
			r.NoError(err)

			rr := httptest.NewRecorder()
			handler.GetUserByID(rr, req)

			r.Equal(tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)

	log, err := logger.New("mock")
	r.NoError(err)

	handler := New(log, mockUserService)

	tests := []struct {
		name           string
		requestBody    string
		mockBehavior   func(user entity.User)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"name": "Test", "email": "test@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().CreateUser(gomock.Any(), user).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"name": "Test", "email": "test@example.com"`,
			mockBehavior:   func(user entity.User) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "user already exists",
			requestBody: `{"name": "test", "email": "test@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().CreateUser(gomock.Any(), user).Return(entity.ErrAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal server error",
			requestBody: `{"name": "test", "email": "test@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().CreateUser(gomock.Any(), user).Return(errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user entity.User
			if tt.name != "invalid request body" {
				err := json.Unmarshal([]byte(tt.requestBody), &user)
				r.NoError(err)
				tt.mockBehavior(user)
			}

			req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.requestBody))
			r.NoError(err)

			rr := httptest.NewRecorder()
			handler.CreateUser(rr, req)

			r.Equal(tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)

	log, err := logger.New("mock")
	r.NoError(err)

	handler := New(log, mockUserService)

	tests := []struct {
		name           string
		requestBody    string
		mockBehavior   func(user entity.User)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"id": "d290f1ee-6c54-4b01-90e6-d701748f0851", "name": "Updated Name", "email": "updated@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().UpdateUser(gomock.Any(), user).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"id": "d290f1ee-6c54-4b01-90e6-d701748f0851", "name": "Updated Name", "email": "updated@example.com"`,
			mockBehavior:   func(user entity.User) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "user not found",
			requestBody: `{"id": "d290f1ee-6c54-4b01-90e6-d701748f0851", "name": "Updated Name", "email": "updated@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().UpdateUser(gomock.Any(), user).Return(entity.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "internal server error",
			requestBody: `{"id": "d290f1ee-6c54-4b01-90e6-d701748f0851", "name": "Updated Name", "email": "updated@example.com"}`,
			mockBehavior: func(user entity.User) {
				mockUserService.EXPECT().UpdateUser(gomock.Any(), user).Return(errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user entity.User
			if tt.name != "invalid request body" {
				err := json.Unmarshal([]byte(tt.requestBody), &user)
				r.NoError(err)
				tt.mockBehavior(user)
			}

			req, err := http.NewRequest(http.MethodPut, "/user", strings.NewReader(tt.requestBody))
			r.NoError(err)

			rr := httptest.NewRecorder()
			handler.UpdateUser(rr, req)

			r.Equal(tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)

	log, err := logger.New("mock")
	r.NoError(err)

	handler := New(log, mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockBehavior   func(userID uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "success",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing id",
			userID:         "",
			mockBehavior:   func(userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid id",
			userID:         "invalid-id",
			mockBehavior:   func(userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().DeleteUser(gomock.Any(), userID).Return(entity.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "internal server error",
			userID: uuid.Must(uuid.NewV4()).String(),
			mockBehavior: func(userID uuid.UUID) {
				mockUserService.EXPECT().DeleteUser(gomock.Any(), userID).Return(errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success" || tt.name == "internal server error" || tt.name == "user not found" {
				userID, err := uuid.FromString(tt.userID)
				r.NoError(err)
				tt.mockBehavior(userID)
			}

			req, err := http.NewRequest(http.MethodDelete, "/user?id="+tt.userID, nil)
			r.NoError(err)

			rr := httptest.NewRecorder()
			handler.DeleteUser(rr, req)

			r.Equal(tt.expectedStatus, rr.Code)
		})
	}
}
