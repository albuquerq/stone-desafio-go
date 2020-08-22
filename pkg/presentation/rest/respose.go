package rest

// Response is the result of a rest api request.
type Response struct {
	Error *Error      `json:"error,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Error is the rest api error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	GoErr   error  `json:"go_err"`
}
