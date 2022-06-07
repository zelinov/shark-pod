package exception

import "net/http"

var httpStatuses = map[Type]int{
	BadRequest:          http.StatusBadRequest,
	Unauthorized:        http.StatusUnauthorized,
	Forbidden:           http.StatusForbidden,
	NotFound:            http.StatusNotFound,
	NotAcceptable:       http.StatusNotAcceptable,
	Conflict:            http.StatusConflict,
	UnprocessableEntity: http.StatusUnprocessableEntity,
	InternalServerError: http.StatusInternalServerError,
}

func HttpResponseEntity(err error) (msg string, status int, code string) {
	status = http.StatusInternalServerError

	var e *Exception

	switch err.(type) {
	case Exception:
		te, _ := err.(Exception)
		e = &te
	case *Exception:
		e, _ = err.(*Exception)
	}

	if e != nil {
		msg = e.Error()
		code = string(e.Code())

		if s, ok := httpStatuses[e.Type()]; ok {
			status = s
		}
	} else {
		msg = err.Error()
		code = string(Unknown)
	}

	return
}
