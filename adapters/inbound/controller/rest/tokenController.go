package rest

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	gateway "sharks/adapters/inbound/controller"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/ports/inbound"
	"sharks/utils/jwt"
)

type tokenController struct {
	service inbound.TokenService
}

func initTokenController(r *httprouter.Router, s inbound.TokenService) {
	c := tokenController{s}

	r.GET("/tokens", c.getAllByOwner)
}

func (c *tokenController) getAllByOwner(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if pk, err := jwt.GetPublicKey(gateway.GetAccessToken(r)); err != nil {
		logger.Log.Error(err.Error())

		gateway.WriteError(rw, r, err)
	} else {
		if tokens, err := c.service.GetAll(&application.TokenFilter{Owner: pk, Creators: nil}); err == nil {
			var imgs []ImageDto

			for _, t := range tokens {
				imgs = append(imgs, ImageDto{t.ImageUrl})
			}

			writeResponse(rw, r, GetAllTokensByOwnerResponse{imgs})
		} else {
			logger.Log.Error(err.Error())

			gateway.WriteError(rw, r, err)
		}
	}
}
