package response

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

// TODO: make error argument not string
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
