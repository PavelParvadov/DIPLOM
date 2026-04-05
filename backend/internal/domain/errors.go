package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrAlreadyMember      = errors.New("already a member of this house")
	ErrValidation         = errors.New("validation error")
	ErrExpiredInviteCode  = errors.New("invite code expired")
)
