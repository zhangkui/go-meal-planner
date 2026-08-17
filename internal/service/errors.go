package service

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrMenuNotReplaceable = errors.New("unconfirmed menu already exists")
	ErrInsufficientStock  = errors.New("insufficient inventory")
	ErrUnitMismatch       = errors.New("ingredient unit mismatch")
)
