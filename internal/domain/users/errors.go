package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrNoRefreshToken = errors.New("refresh token not exist")
var ErrInvalidRefreshToken = errors.New("refresh token is not valid")
var ErrInvalidAccessToken = errors.New("access token is not valid")

type NoUserError struct {
	GUID uuid.UUID
}

func NewNoUserError(guid uuid.UUID) error {
	return NoUserError{GUID: guid}
}

func (e NoUserError) Error() string {
	return fmt.Sprintf("user with uuid: %s does not exist", e.GUID)
}
