package utils

type (
	WSApiError struct {
		Code  int
		Error string
	}
	WSApiRequest struct {
		Func    string
		Payload interface{}
	}
	WSApiReply struct {
		Ok      bool
		Error   *WSApiError
		Payload interface{}
	}
)
