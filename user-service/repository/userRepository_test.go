package repository

import (
	"testing"

	"github.com/dharmasatrya/goodkarma/user-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUserSupporter(request entity.CreateUserSupporterRequest) (*entity.User, error) {
	args := m.Called(request)

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserCoordinator(request entity.CreateUserCoordinatorRequest) (*entity.User, error) {
	args := m.Called(request)

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Login(request entity.LoginRequest) (*entity.User, error) {
	args := m.Called(request)

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id string) (*entity.DetailUser, error) {
	args := m.Called(id)

	return args.Get(0).(*entity.DetailUser), args.Error(1)
}

func (m *MockUserRepository) UpdateProfile(request entity.UpdateProfileRequest) (*entity.DetailUser, error) {
	args := m.Called(request)

	return args.Get(0).(*entity.DetailUser), args.Error(1)
}

func TestCreateUserSupporter_Success(t *testing.T) {
	mockUserRepository := new(MockUserRepository)

	user := entity.CreateUserSupporterRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "password",
		Role:     "supporter",
		FullName: "Test User",
		Address:  "Test Address",
		Phone:    "081234567890",
		Photo:    "test.jpg",
	}

	mockUserRepository.On("CreateUserSupporter", user).Return(user, nil)

	userRepository := &userRepository{userRepository: mockUserRepository}

	result, err := userRepository.CreateUserSupporter(user)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	mockUserRepository.AssertExpectations(t)
}

func TestCreateUserSupporter_Failed(t *testing.T) {
	mockUserRepository := new(MockUserRepository)

	user := entity.CreateUserSupporterRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "password",
		Role:     "supporter",
		FullName: "Test User",
		Address:  "Test Address",
		Phone:    "081234567890",
		Photo:    "test.jpg",
	}

	mockUserRepository.On("CreateUserSupporter", user).Return(nil, assert.AnError)

	userRepository := &userRepository{userRepository: mockUserRepository}

	result, err := userRepository.CreateUserSupporter(user)

	assert.Nil(t, result)
	assert.NotNil(t, err)

	mockUserRepository.AssertExpectations(t)
}
