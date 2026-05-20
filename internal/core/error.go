package core

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCodeBadRequest           = errors.New("ERR_BAD_REQUEST")
	ErrCodeInternalError        = errors.New("ERR_INTERNAL_ERROR")
	ErrSessionExpired           = errors.New("ERR_SESSION_EXPIRED")
	ErrMessageNotFound          = errors.New("ERR_MESSAGE_NOT_FOUND")
	ErrCannotRetryActiveMessage = errors.New("ERR_CANNOT_RETRY_ACTIVE_MESSAGE")
)

// CannotRetryMessageError indicates a retry was requested for a message that is not failed.
type CannotRetryMessageError struct {
	CurrentStatus MessageStatus
}

func NewCannotRetryMessageError(status MessageStatus) *CannotRetryMessageError {
	return &CannotRetryMessageError{CurrentStatus: status}
}

func (e *CannotRetryMessageError) Error() string {
	return fmt.Sprintf(
		"Cannot retry message: current status is '%s'. Only '%s' messages can be retried.",
		strings.ToUpper(string(e.CurrentStatus)),
		strings.ToUpper(string(MessageStatusFailed)),
	)
}

func (e *CannotRetryMessageError) Is(target error) bool {
	return target == ErrCannotRetryActiveMessage
}
