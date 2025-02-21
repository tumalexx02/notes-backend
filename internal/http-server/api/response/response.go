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

func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
