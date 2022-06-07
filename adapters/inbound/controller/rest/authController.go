package rest

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	gateway "sharks/adapters/inbound/controller"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/application/exception"
	"sharks/ports/inbound"
	"sharks/utils/jwt"
)

type authController struct {
	authService inbound.AuthService
}

func initAuthController(r *httprouter.Router, as inbound.AuthService) {
	c := authController{as}

	r.POST("/login", c.login)
	r.POST("/refresh", c.refresh)
	r.GET("/nonce/:pk", c.nonce)
	r.POST("/check", c.check)
	r.POST("/logout", c.logout)
}

func (c *authController) login(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error(err.Error())

		e := exception.FromString("Wrong credentials", exception.BadRequest, exception.InvalidToken)
		gateway.WriteError(rw, r, e)

		return
	}

	cred := application.Credentials{
		Signature: req.Signature,
		PublicKey: req.PublicKey,
	}

	if token, e := c.authService.Login(cred); e != nil {
		logger.Log.Error(e.Error())

		gateway.WriteError(rw, r, e)
	} else {
		writeResponse(rw, r, LoginResponse{
			Access:    token.Access.Key,
			Refresh:   token.Refresh.Key,
			PublicKey: token.Access.PublicKey,
			ExpiredAt: token.Access.Expired,
		})
	}
}

func (c *authController) refresh(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error(err.Error())

		e := exception.FromString("Wrong credentials", exception.BadRequest, exception.InvalidToken)
		gateway.WriteError(rw, r, e)

		return
	}

	if token, e := c.authService.Refresh(req.Refresh); e != nil {
		logger.Log.Error(e.Error())

		gateway.WriteError(rw, r, e)
	} else {
		writeResponse(rw, r, RefreshResponse{
			Access:    token.Access.Key,
			Refresh:   token.Refresh.Key,
			PublicKey: token.Access.PublicKey,
			ExpiredAt: token.Access.Expired,
		})
	}
}

func (c *authController) nonce(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pk := ps.ByName("pk")
	nonce := c.authService.Nonce(pk)

	writeResponse(rw, r, NonceResponse{nonce.Nonce.String()})
}

func (c *authController) check(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if pk, err := jwt.GetPublicKey(gateway.GetAccessToken(r)); err != nil {
		logger.Log.Error(err.Error())

		gateway.WriteError(rw, r, err)
	} else {
		writeResponse(rw, r, CheckResponse{pk})
	}
}

func (c *authController) logout(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if e := c.authService.Logout(gateway.GetAccessToken(r)); e != nil {
		gateway.WriteError(rw, r, e)

		return
	}

	writeResponse(rw, r, nil)
}
