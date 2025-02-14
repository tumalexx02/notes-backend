package response

const (
	StatusOk                  = "OK"
	StatusError               = "ERROR"
	StatusUnauthorized        = "UNAUTHORIZED"
	StatusExpiredAccessToken  = "REVOKED_ACCESS_TOKEN"
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

func Unauthorized(msg string) Response {
	return Response{
		Status: StatusUnauthorized,
		Error:  msg,
	}
}

func ExpiredAccessToken() Response {
	return Response{
		Status: StatusExpiredAccessToken,
	}
}

func RevokedRefreshToken() Response {
	return Response{
		Status: StatusRevokedRefreshToken,
	}
}
