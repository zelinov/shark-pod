package exception

const (
	InternalServerError Type = iota
	Unauthorized
	Forbidden
	NotFound
	UnprocessableEntity
	BadRequest
	Conflict
	NotAcceptable
)

type Type uint16

type Exception struct {
	message string
	eType   Type
	code    Code
}

func (e Exception) Error() string {
	return e.message
}

func (e Exception) Type() Type {
	return e.eType
}

func (e Exception) Code() Code {
	return e.code
}

func FromError(err error, eType Type, code Code) *Exception {
	return &Exception{
		message: err.Error(),
		eType:   eType,
		code:    code,
	}
}

func FromException(err error) *Exception {
	switch err.(type) {
	case Exception:
		e := err.(Exception)

		return &e

	case *Exception:
		return err.(*Exception)

	default:
		return FromError(err, InternalServerError, Unknown)
	}
}

func FromString(message string, eType Type, code Code) *Exception {
	return &Exception{
		message: message,
		eType:   eType,
		code:    code,
	}
}
