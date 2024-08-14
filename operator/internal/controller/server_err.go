package controller

// serverErr struct is used to parse error response from our news-aggregator server.
//
// When go-gator server encounters errors with something, it returns them in json format with
// one error field, which indicates the actual error message.
// This will be used to print
type serverErr struct {
	ErrorMsg string `json:"error"`
}

// Error function is used to implement error interface
//
// It returns the error with which this struct was initialized.
func (e *serverErr) Error() string {
	return e.ErrorMsg
}
