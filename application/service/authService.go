package service

import (
	"fmt"
	"github.com/google/uuid"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/application/exception"
	"sharks/config"
	"sharks/ports/inbound"
	"sharks/ports/outbound"
	"sharks/utils/jwt"
)

type JwtAuthService struct {
	jwtRepository   outbound.JwtRepository
	nonceRepository outbound.NonceRepository
	solanaService   inbound.SolanaService
	tokenService    inbound.TokenService
}

func NewJwtAuthService(
	jr outbound.JwtRepository,
	nr outbound.NonceRepository,
	ss inbound.SolanaService,
	ts inbound.TokenService,
) *JwtAuthService {
	return &JwtAuthService{jr, nr, ss, ts}
}

var creatorAccounts = config.GetConfig().CreatorAccounts
var creatorNftCount = config.GetConfig().CreatorNftCount

func (s *JwtAuthService) Login(cred application.Credentials) (res *application.JwtToken, e *exception.Exception) {
	pk := cred.PublicKey

	if nonceVar, err := s.nonceRepository.Get(pk); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidSignature)

		return
	} else if e = s.solanaService.VerifySignature(cred.Signature, pk, nonceVar.Nonce.String()); e != nil {
		return
	}

	var tokens []*application.TokenMetadata

	if tokens, e = s.tokenService.GetAll(&application.TokenFilter{Owner: pk, Creators: &creatorAccounts}); e != nil {
		return
	} else if len(tokens) < creatorNftCount {
		e = exception.FromString(
			fmt.Sprintf("Less than %d sharks tokens", creatorNftCount),
			exception.Unauthorized,
			exception.InsufficientNft,
		)

		return
	}

	if t, err := jwt.GetJwtToken(pk); err == nil {
		if err := s.jwtRepository.DeleteAllByPublicKey(pk); err != nil {
			logger.Log.Error(err.Error())
		}

		if e = s.createAuth(t); e == nil {
			res = t
		}
	} else {
		logger.Log.Error(err.Error())

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
	}

	return
}

func (s *JwtAuthService) Verify(strTkn string) (e *exception.Exception) {
	var tokMeta *application.Token

	if strTkn == "" {
		logger.Log.Error("Empty access token")

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)

		return
	}

	if tokMeta, e = verifyAccessToken(strTkn); e == nil {
		if t, err := s.jwtRepository.Get(tokMeta.Id); err == nil {
			if t.PublicKey != tokMeta.PublicKey {
				e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
			}
		} else {
			logger.Log.Error(err.Error())

			e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
		}
	}

	return
}

func (s *JwtAuthService) Refresh(refresh string) (t *application.JwtToken, e *exception.Exception) {
	var pk string

	if pk, e = s.verifyRefresh(refresh); e != nil {
		return
	}

	if err := s.jwtRepository.DeleteAllByPublicKey(pk); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)

		return
	}

	var err error

	if t, err = jwt.GetJwtToken(pk); e == nil {
		e = s.createAuth(t)
	} else {
		logger.Log.Error(err.Error())

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
	}

	return
}

func (s *JwtAuthService) Logout(strTkn string) (e *exception.Exception) {
	var tokMeta *application.Token

	if tokMeta, e = verifyAccessToken(strTkn); e != nil {
		return
	}

	if err := s.jwtRepository.DeleteAllByPublicKey(tokMeta.PublicKey); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
	}

	return
}

func (s *JwtAuthService) Nonce(pk string) *application.Nonce {
	if nonce, e := s.nonceRepository.Get(pk); e == nil {
		return nonce
	}

	nonce := &application.Nonce{
		PublicKey: pk,
		Nonce:     uuid.New(),
	}

	if err := s.nonceRepository.Save(nonce); err != nil {
		logger.Log.Error(err.Error())
	}

	return nonce
}

func verifyAccessToken(access string) (t *application.Token, e *exception.Exception) {
	if claims, err := jwt.ParseAccess(access); err == nil {
		if atId, ok := claims["access_id"].(string); ok {
			if pk, ok := claims["public_key"].(string); ok {
				t = &application.Token{
					Id:        uuid.MustParse(atId),
					Key:       access,
					PublicKey: pk,
				}
			} else {
				e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
			}
		} else {
			e = exception.FromString("Token is invalid", exception.Unauthorized, exception.InvalidToken)
		}
	} else {
		logger.Log.Error(err.Error())

		e = exception.FromString("Token is invalid", exception.Unauthorized, exception.InvalidToken)
	}

	return
}

func (s *JwtAuthService) createAuth(t *application.JwtToken) (e *exception.Exception) {
	if err := s.jwtRepository.Save(&t.Access); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("Something went wrong", exception.Unauthorized, exception.Unknown)
	} else if err := s.jwtRepository.Save(&t.Refresh); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("Something went wrong", exception.Unauthorized, exception.Unknown)
	}

	return
}

func (s *JwtAuthService) verifyRefresh(refresh string) (pk string, e *exception.Exception) {
	if claims, err := jwt.ParseRefresh(refresh); err == nil {
		if rId, ok := claims["refresh_id"].(string); ok {
			if t, err := s.jwtRepository.Get(uuid.MustParse(rId)); err == nil && t != nil {
				if pk, ok = claims["public_key"].(string); !ok {
					e = exception.FromString("Token is invalid", exception.Unauthorized, exception.InvalidToken)
				}
			} else {
				if err != nil {
					logger.Log.Error(err.Error())
				}

				e = exception.FromString("Unauthorized", exception.Unauthorized, exception.InvalidToken)
			}
		} else {
			e = exception.FromString("Token is invalid", exception.Unauthorized, exception.InvalidToken)
		}
	} else {
		logger.Log.Error(err.Error())

		e = exception.FromString("Token is invalid", exception.Unauthorized, exception.InvalidToken)
	}

	return
}
