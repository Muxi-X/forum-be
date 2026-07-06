package errno

import (
	"go-micro.dev/v4/errors"
)

func ServerErr(errno *Errno, cause string) error {
	return &errors.Error{
		Code:   int32(errno.Code),
		Detail: cause,
		Status: errno.Message,
	}
}

func NotFoundErr(errno *Errno, cause string) error {
	return &errors.Error{
		Code:   404,
		Detail: cause,
		Status: errno.Message,
	}
}

func (e Errno) Error() string {
	return e.Message
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	case *errors.Error:
		if typed.Detail == "" {
			return int(typed.Code), typed.Status
		}
		return int(typed.Code), typed.Status + " : " + typed.Detail
	default:
		return InternalServerError.Code, err.Error()
	}
}

// type ErrDetail struct {
// 	Errno
// 	Cause string `json:"cause"`
// }

// func ParseDetail(detail string) (*ErrDetail, error) {
// 	d := &ErrDetail{}
// 	err := json.Unmarshal([]byte(detail), d)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return d, nil
// }

// // Err represents an error
// type Err struct {
// 	Code    int
// 	Message string
// 	Err     error
// }

// func New(errno *Errno, err error) *Err {
// 	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
// }
//
// func (err *Err) Add(message string) error {
// 	err.Message += " " + message
// 	return err
// }
//
// func (err *Err) Addf(format string, args ...interface{}) error {
// 	err.Message += " " + fmt.Sprintf(format, args...)
// 	return err
// }
//
// func (err *Err) Error() string {
// 	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
// }
