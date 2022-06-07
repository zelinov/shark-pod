package controller

import (
	"go.uber.org/zap"
	"net/http"
	"sharks/adapters/outbound/logger"
	"sharks/application/exception"
	"sharks/ports/inbound"
	"sharks/utils/url"
	"strings"
)

const authorizationPrefix = "Bearer "

var filter = url.NewUrlFilter([]string{
	"api/v1/login",
	"api/v1/refresh",
	"api/v1/nonce",
})

func NewHandler(h http.Handler, as inbound.AuthService) http.Handler {
	return loggerMiddleware(securityMiddleware(h, as))
}

func securityMiddleware(h http.Handler, as inbound.AuthService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if ignore, err := filter.IsIgnore(r.URL.Path); !ignore {
			if e := as.Verify(GetAccessToken(r)); e != nil {
				WriteError(rw, r, exception.FromString(e.Error(), exception.Unauthorized, e.Code()))

				return
			}
		} else if err != nil {
			logger.Log.Error(err.Error())

			WriteError(rw, r, exception.FromString("Something went wrong", exception.Unauthorized, exception.Unknown))

			return
		}

		h.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Log.Info(
			"Request",
			zap.String("method", r.Method),
			zap.String("host", r.Host),
			zap.String("uri", r.RequestURI),
		)

		h.ServeHTTP(rw, r)
	})
}

func GetAccessToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), authorizationPrefix)
}

func WriteError(rw http.ResponseWriter, r *http.Request, err error) {
	message, status, code := exception.HttpResponseEntity(err)

	rw.Header().Set("Error-Code", code)
	http.Error(rw, message, status)

	logger.Log.Info(
		"Response error",
		zap.String("method", r.Method),
		zap.String("host", r.Host),
		zap.String("uri", r.RequestURI),
		zap.Int("status", status),
		zap.String("code", code),
		zap.String("message", message),
	)
}
