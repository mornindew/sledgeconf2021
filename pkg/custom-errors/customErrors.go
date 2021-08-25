package customerrors

import "strconv"

//Generic error for handling precondition errors
type PreconditionError struct {
	Msg string
}

func (e PreconditionError) Error() string {
	return e.Msg
}

//Generic error for anythign with invalid data
type InvalidData struct {
	Msg               string
	InternalErrorCode int
}

func (e InvalidData) Error() string {
	return e.Msg + " Error Code: " + strconv.Itoa(e.InternalErrorCode)
}

//Bad Format Error
type BadFormat struct {
	Msg string
}

func (e BadFormat) Error() string {
	return e.Msg
}

//InternalServerError - since it is generic it has an internal server code on it
type InternalServerError struct {
	Msg               string
	InternalErrorCode int
}

func (e InternalServerError) Error() string {
	return e.Msg + " Error Code: " + strconv.Itoa(e.InternalErrorCode)
}

//Client Connection ERror
type ClientConstructionError struct {
	Msg string
}

func (e ClientConstructionError) Error() string {
	return e.Msg
}
