package application

import (
	"net/http"

	"github.com/pkg/errors"
)

func HandleError(err error) error {
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

type ErrHTTP interface {
	Code() int
	Error() string
}

type ErrCustom struct {
	code    int
	message string
}

func NewErrCustom(code int, msg string) ErrCustom {
	return ErrCustom{code: code, message: msg}
}

func (self ErrCustom) Error() string {
	return self.message
}

func (self ErrCustom) Code() int {
	return self.code
}

type ErrConflict struct {
	message string
}

func NewErrConflict(msg string) ErrConflict {
	return ErrConflict{message: msg}
}

func (self ErrConflict) Error() string {
	return self.message
}

func (self ErrConflict) Code() int {
	return http.StatusConflict
}

type ErrNotFound struct {
	message string
}

func NewErrNotFound(name string) ErrNotFound {
	return ErrNotFound{message: name + " not found"}
}

func (self ErrNotFound) Error() string {
	return self.message
}

func (self ErrNotFound) Code() int {
	return http.StatusNotFound
}

type ErrValidation struct {
	message string
}

func NewErrValidation(msg string) ErrValidation {
	return ErrValidation{message: msg}
}

func (self ErrValidation) Error() string {
	return self.message
}

func (self ErrValidation) Code() int {
	return http.StatusBadRequest
}

type ErrUnauthorized struct{}

func NewErrUnauthorized() ErrUnauthorized { return ErrUnauthorized{} }

func (self ErrUnauthorized) Error() string {
	return "Unauthorized"
}

func (self ErrUnauthorized) Code() int {
	return http.StatusUnauthorized
}

type ErrForbidden struct{}

func NewErrForbidden() ErrForbidden { return ErrForbidden{} }

func (self ErrForbidden) Error() string {
	return "Forbiden"
}

func (self ErrForbidden) Code() int {
	return http.StatusForbidden
}
