package webhttp

type ResponseError struct {
	Status       string `json:"status"`
	StatusCode   uint16 `json:"statusCode"`
	StatusText   string `json:"statusText"`
	ErrorMessage string `json:"error"`
}

type ResponseErrors struct {
	Status        string   `json:"status"`
	StatusCode    uint16   `json:"statusCode"`
	StatusText    string   `json:"statusText"`
	ErrorMessages []string `json:"errors"`
}

type ResponseSuccess struct {
	Status     string      `json:"status"`
	StatusCode uint16      `json:"statusCode"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
}

func NewOk(data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Status:     "SUCCESS",
		StatusCode: 200,
		StatusText: "OK",
		Data:       data,
	}
}

func NewBadRequestValidation(errorMessages []string) ResponseErrors {
	return ResponseErrors{
		Status:        "ERROR",
		StatusCode:    400,
		StatusText:    "BAD_REQUEST",
		ErrorMessages: errorMessages,
	}
}

func NewBadRequest(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   400,
		StatusText:   "BAD_REQUEST",
		ErrorMessage: errorMessage,
	}

}
func NewUnauthorizedRequest(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   401,
		StatusText:   "UNAUTHORIZED",
		ErrorMessage: errorMessage,
	}
}

func NewForbiddenRequest(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   401,
		StatusText:   "FORBIDDEN",
		ErrorMessage: errorMessage,
	}
}

func NewNotFound(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   404,
		StatusText:   "NOT_FOUND",
		ErrorMessage: errorMessage,
	}
}

func NewConflict(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   409,
		StatusText:   "CONFLICT",
		ErrorMessage: errorMessage,
	}
}

func NewInternalServerError(errorMessage string) ResponseError {
	return ResponseError{
		Status:       "ERROR",
		StatusCode:   500,
		StatusText:   "INTERNAL_SERVER_ERROR",
		ErrorMessage: errorMessage,
	}
}
