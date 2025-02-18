package response

const (
	StatusOk                  = "OK"
	StatusError               = "ERROR"
	StatusRevokedRefreshToken = "REVOKED_REFRESH_TOKEN"
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

func RevokedRefreshToken() Response {
	return Response{
		Status: StatusRevokedRefreshToken,
	}
}
