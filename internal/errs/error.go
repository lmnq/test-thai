package errs

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUniqueConstraint  = errors.New("unique constraint error")
	UniqueConstraintCode = "23505"

	// status code error messages
	StatusBadRequestMessage          = "bad request"
	StatusNotFoundMessage            = "not found"
	StatusInternalServerErrorMessage = "internal server error"
)
