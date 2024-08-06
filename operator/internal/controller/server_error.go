package controller

// serverError struct is used to represent the error message of the server.
type serverError struct {
	ErrMsg string `json:"error"`
}

// Error function is used to implement error interface.
// It returns the error message of the serverError struct.
func (e *serverError) Error() string {
	return e.ErrMsg
}
