package rest

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	gateway "sharks/adapters/inbound/controller"
	"sharks/application/exception"
	"sharks/ports/inbound"
)

func NewHttpHandler(as inbound.AuthService, ts inbound.TokenService) http.Handler {
	r := httprouter.New()

	initAuthController(r, as)
	initTokenController(r, ts)

	return http.StripPrefix("/api/v1", r)
}

func writeResponse(rw http.ResponseWriter, r *http.Request, t interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if t != nil {
		if err := json.NewEncoder(rw).Encode(t); err != nil {
			gateway.WriteError(rw, r, exception.FromError(err, exception.InternalServerError, exception.Unknown))
		}
	}
}
