package customerrors

import "strconv"

//PreconditionError - This is used if any mandatory data needed for a function to complete isn't available.
type PreconditionError struct {
	Msg string
}

func (e PreconditionError) Error() string {
	return e.Msg
}

//InvalidData - Error that will occur if the data is present but isn't correct
type InvalidData struct {
	Msg               string
	InternalErrorCode int
}

func (e InvalidData) Error() string {
	return e.Msg + " Error Code: " + strconv.Itoa(e.InternalErrorCode)
}

//Bad Format Error - used when the format of the data is incorrect
type BadFormat struct {
	Msg string
}

func (e BadFormat) Error() string {
	return e.Msg
}

//Bad Request - This is used when the actual request (typically HTTP) isn't formatted correctly
type BadRequest struct {
	Msg string
}

func (e BadRequest) Error() string {
	return e.Msg
}

//Not Found Error - used if an entity isn't available
type NotFoundError struct {
	Msg string
}

func (e NotFoundError) Error() string {
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

//Client Connection Error - Custom error for when we are unable to create a client connection
type ClientConstructionError struct {
	Msg string
}

func (e ClientConstructionError) Error() string {
	return e.Msg
}

/*
HTTPError - Generic Error for unhandled HTTP errors
*/
type HTTPError struct {
	Msg  string
	Code int
}

func (e HTTPError) Error() string {
	return e.Msg + " Error Code: " + strconv.Itoa(e.Code)
}
