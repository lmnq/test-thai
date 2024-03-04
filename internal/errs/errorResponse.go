package errs

type Error struct {
	Err     error  `json:"err,omitempty"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// func (e Error) Error() string {
// 	return e.Err.Error()
// }

func NilError() Error {
	return Error{
		Err:     nil,
		Code:    0,
		Message: "",
	}
}

func (e Error) IsErr() bool {
	return e.Err != nil
}