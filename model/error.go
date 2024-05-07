package model

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(`{"code": %d, "message": "%s"}`, e.Code, e.String())
}

func (e *Error) String() string {
	return e.Message
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"code": %d, "message": "%s"}`, e.Code, e.String())), nil
}
