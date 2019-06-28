package session

import "errors"

var(
	ErrorSessionValueNotFound = errors.New("session value not found")
	ErrorSessionNotBeNil = errors.New("session  not be nil")
	ErrorSessionNotFound = errors.New("session  not found")
	ErrorSessionIdNotBeEmpty = errors.New("session id not be empty")
	ErrorSessionTimeoutNotBeNegative = errors.New("session timeout  not be negative")
)
