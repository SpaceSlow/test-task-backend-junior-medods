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

func TestUserService_RefreshTokens_SuccessRefreshing(t *testing.T) {
	email := "example@example.org"
	tokenLifetime := time.Hour
	secretKey := "test-key"

	userGUID, err := uuid.NewRandom()
	require.NoError(t, err, "Не удалось сгенерировать рандомный GUID пользователя")

	ip := net.ParseIP("172.0.19.1")
	require.NotNil(t, ip, "IP-адрес пустой")

	access, err := users.GenerateAccessToken(email, ip, tokenLifetime, secretKey)
	require.NoError(t, err, "Не удалось сгенерировать Access токен")

	refresh, err := users.NewRefreshToken()
	require.NoError(t, err, "Не удалось сгенерировать Refresh токен")
	refreshHash, err := refresh.GenerateHash()
	require.NoError(t, err, "Не удалось вычислить хэш для Refresh токена")
	user := users.NewUser(userGUID, email, string(refreshHash))

	mockRepo := mocks.NewRepository(t)
	mockRepo.On("FetchEmailByUUID", userGUID).Return(email, nil)
	mockRepo.On("CreateRefreshToken", userGUID, mock.Anything).Return(nil)
	mockRepo.On("FetchUserByEmail", email).Return(user, nil)

	mockNotifier := mocks.NewNotifierService(t)

	mockCfg := mocks.NewConfig(t)
	mockCfg.On("TokenLifetime").Return(tokenLifetime)
	mockCfg.On("SecretKey").Return(secretKey)

	s := NewUserService(mockRepo, mockNotifier, mockCfg)

	time.Sleep(time.Second) // for generating access tokens with different expired_at field value
	gotAccess, gotRefresh, err := s.RefreshTokens(access, refresh, ip)
	require.NoError(t, err, "Возникла ошибка при получении Access и Refresh токенов")
	require.NotNil(t, gotAccess, "При получении токенов Access токен не был получен")
	require.NotEqual(t, access, gotAccess, "После обновления токенов Access токены совпадают")
	require.NotNil(t, gotRefresh, "При получении токенов Refresh токен не был получен")
	require.NotEqual(t, refresh, gotRefresh, "После обновления токенов Refresh токены совпадают")

	assert.NotEmpty(t, gotAccess.JWT(), "JWT Access токена оказался пустым")
	err = gotAccess.Parse(secretKey)
	require.NoError(t, err, "Возникла ошибка при извлечении данных из Access токена")

	assert.Equal(t, ip, gotAccess.IP(), "Не совпадает ip в Access токене")
	assert.Equal(t, email, gotAccess.Email(), "Не совпадает email в Access токене")
}

func TestUserService_RefreshTokens_SuccessRefreshingWithNewIP(t *testing.T) {
	email := "example@example.org"
	tokenLifetime := time.Hour
	secretKey := "test-key"

	userGUID, err := uuid.NewRandom()
	require.NoError(t, err, "Не удалось сгенерировать рандомный GUID пользователя")

	oldIP := net.ParseIP("172.0.19.1")
	require.NotNil(t, oldIP, "IP-адрес пустой")

	newIP := net.ParseIP("172.0.20.1")
	require.NotNil(t, newIP, "IP-адрес пустой")

	access, err := users.GenerateAccessToken(email, oldIP, tokenLifetime, secretKey)
	require.NoError(t, err, "Не удалось сгенерировать Access токен")

	refresh, err := users.NewRefreshToken()
	require.NoError(t, err, "Не удалось сгенерировать Refresh токен")
	refreshHash, err := refresh.GenerateHash()
	require.NoError(t, err, "Не удалось вычислить хэш для Refresh токена")
	user := users.NewUser(userGUID, email, string(refreshHash))

	mockRepo := mocks.NewRepository(t)
	mockRepo.On("FetchEmailByUUID", userGUID).Return(email, nil)
	mockRepo.On("CreateRefreshToken", userGUID, mock.Anything).Return(nil)
	mockRepo.On("FetchUserByEmail", email).Return(user, nil)

	mockNotifier := mocks.NewNotifierService(t)
	mockNotifier.On("SendSuspiciousActivityMail", email, newIP).Return(nil)

	mockCfg := mocks.NewConfig(t)
	mockCfg.On("TokenLifetime").Return(tokenLifetime)
	mockCfg.On("SecretKey").Return(secretKey)

	time.Sleep(time.Second) // for generating access tokens with different expired_at field value
	s := NewUserService(mockRepo, mockNotifier, mockCfg)

	gotAccess, gotRefresh, err := s.RefreshTokens(access, refresh, newIP)
	require.NoError(t, err, "Возникла ошибка при получении Access и Refresh токенов")
	require.NotNil(t, gotAccess, "При получении токенов Access токен не был получен")
	require.NotEqual(t, access, gotAccess, "После обновления токенов Access токены совпадают")
	require.NotNil(t, gotRefresh, "При получении токенов Refresh токен не был получен")
	require.NotEqual(t, refresh, gotRefresh, "После обновления токенов Refresh токены совпадают")

	assert.NotEmpty(t, gotAccess.JWT(), "JWT Access токена оказался пустым")
	err = gotAccess.Parse(secretKey)
	require.NoError(t, err, "Возникла ошибка при извлечении данных из Access токена")

	assert.Equal(t, newIP, gotAccess.IP(), "Не совпадает ip в Access токене")
	assert.Equal(t, email, gotAccess.Email(), "Не совпадает email в Access токене")
}

func TestUserService_RefreshTokens_InvalidRefreshToken(t *testing.T) {
	email := "example@example.org"
	tokenLifetime := time.Hour
	secretKey := "test-key"

	ip := net.ParseIP("172.0.19.1")
	require.NotNil(t, ip, "IP-адрес пустой")

	access, err := users.GenerateAccessToken(email, ip, tokenLifetime, secretKey)
	require.NoError(t, err, "Не удалось сгенерировать Access токен")

	refresh, err := users.NewRefreshToken()
	require.NoError(t, err, "Не удалось сгенерировать Refresh токен")

	mockRepo := mocks.NewRepository(t)
	mockRepo.On("FetchUserByEmail", email).Return(nil, users.ErrInvalidRefreshToken)

	mockNotifier := mocks.NewNotifierService(t)

	mockCfg := mocks.NewConfig(t)
	mockCfg.On("SecretKey").Return(secretKey)

	s := NewUserService(mockRepo, mockNotifier, mockCfg)

	invalidRefresh := append(*refresh, '0')
	gotAccess, gotRefresh, err := s.RefreshTokens(access, &invalidRefresh, ip)
	require.Error(t, err, "Не возникло ошибки при получении Access и Refresh токенов для несуществующего пользователя")
	require.Nil(t, gotAccess, "Был получен Access токен для несуществующего пользователя")
	require.Nil(t, gotRefresh, "Был получен Refresh токен для несуществующего пользователя")
}

func TestUserService_RefreshTokens_ExpiredAccessToken(t *testing.T) {
	email := "example@example.org"
	tokenLifetime := time.Duration(0)
	secretKey := "test-key"

	ip := net.ParseIP("172.0.19.1")
	require.NotNil(t, ip, "IP-адрес пустой")

	access, err := users.GenerateAccessToken(email, ip, tokenLifetime, secretKey)
	require.NoError(t, err, "Не удалось сгенерировать Access токен")

	refresh, err := users.NewRefreshToken()
	require.NoError(t, err, "Не удалось сгенерировать Refresh токен")

	mockRepo := mocks.NewRepository(t)

	mockNotifier := mocks.NewNotifierService(t)

	mockCfg := mocks.NewConfig(t)
	mockCfg.On("SecretKey").Return(secretKey)

	s := NewUserService(mockRepo, mockNotifier, mockCfg)

	gotAccess, gotRefresh, err := s.RefreshTokens(access, refresh, ip)
	require.Error(t, err, "Не возникло ошибки при получении Access и Refresh токенов для несуществующего пользователя")
	require.Nil(t, gotAccess, "Был получен Access токен для несуществующего пользователя")
	require.Nil(t, gotRefresh, "Был получен Refresh токен для несуществующего пользователя")
}

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
