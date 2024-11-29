package users

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/service/users/mocks"
)

func TestUserService_Tokens_UserExist(t *testing.T) {
	email := "example@example.org"
	tokenLifetime := time.Hour
	secretKey := "test-key"

	userGUID, err := uuid.NewRandom()
	require.NoError(t, err, "Не удалось сгенерировать рандомный GUID пользователя")

	ip := net.ParseIP("172.0.19.1")
	require.NotNil(t, ip, "IP-адрес пустой")

	mockRepo := mocks.NewRepository(t)
	mockRepo.On("FetchEmailByUUID", userGUID).Return(email, nil)
	mockRepo.On("CreateRefreshToken", userGUID, mock.Anything).Return(nil)

	mockCfg := mocks.NewConfig(t)
	mockCfg.On("TokenLifetime").Return(tokenLifetime)
	mockCfg.On("SecretKey").Return(secretKey)

	s := NewUserService(mockRepo, nil, mockCfg)
	access, refresh, err := s.Tokens(userGUID, ip)
	require.NoError(t, err, "Возникла ошибка при получении Access и Refresh токенов")
	require.NotNil(t, access, "При получении токенов Access токен не был получен")
	require.NotNil(t, refresh, "При получении токенов Refresh токен не был получен")

	assert.NotEmpty(t, access.JWT(), "JWT Access токена оказался пустым")
	err = access.Parse(secretKey)
	require.NoError(t, err, "Возникла ошибка при извлечении данных из Access токена")

	assert.Equal(t, ip, access.IP(), "Не совпадает ip в Access токене")
	assert.Equal(t, email, access.Email(), "Не совпадает email в Access токене")
}

func TestUserService_Tokens_UserNotExist(t *testing.T) {
	userGUID, err := uuid.NewRandom()
	require.NoError(t, err, "Не удалось сгенерировать рандомный GUID пользователя")

	ip := net.ParseIP("172.0.19.1")
	require.NotNil(t, ip, "IP-адрес пустой")

	mockRepo := mocks.NewRepository(t)
	mockRepo.On("FetchEmailByUUID", userGUID).Return("", users.NewNoUserError(userGUID))

	s := NewUserService(mockRepo, nil, nil)
	access, refresh, err := s.Tokens(userGUID, ip)
	require.Error(t, err, "Не возникло ошибки при получении Access и Refresh токенов для несуществующего пользователя")
	require.Nil(t, access, "Был получен Access токен для несуществующего пользователя")
	require.Nil(t, refresh, "Был получен Refresh токен для несуществующего пользователя")
}
