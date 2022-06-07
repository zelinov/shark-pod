package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"sharks/application"
	"sharks/config"
	"time"
)

var accessSecret = config.GetConfig().AccessSecret
var refreshSecret = config.GetConfig().RefreshSecret
var accessExpiration = config.GetConfig().AccessExpiration
var refreshExpiration = config.GetConfig().RefreshExpiration

func GetJwtToken(pk string) (token *application.JwtToken, err error) {
	atExp := time.Now().Add(time.Duration(accessExpiration * int(time.Minute)))
	rtExp := time.Now().Add(time.Duration(refreshExpiration * int(time.Minute)))
	atId := uuid.New()
	rtId := uuid.New()

	atClaims := jwt.MapClaims{
		"authorized": true,
		"access_id":  atId,
		"public_key": pk,
		"exp":        atExp,
	}

	rtClaims := jwt.MapClaims{
		"refresh_id": rtId,
		"public_key": pk,
		"exp":        rtExp,
	}

	var at string
	var rt string

	if at, err = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(accessSecret); err == nil {
		if rt, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString(refreshSecret); err == nil {
			token = &application.JwtToken{
				Access: application.Token{
					Id:        atId,
					Key:       at,
					PublicKey: pk,
					Expired:   atExp,
				},
				Refresh: application.Token{
					Id:        rtId,
					Key:       rt,
					PublicKey: pk,
					Expired:   rtExp,
				},
			}
		}
	}

	return
}

func ParseAccess(access string) (claims jwt.MapClaims, err error) {
	return Parse(access, accessSecret)
}

func ParseRefresh(refresh string) (claims jwt.MapClaims, err error) {
	return Parse(refresh, refreshSecret)
}

func Parse(token string, secret []byte) (claims jwt.MapClaims, err error) {
	var t *jwt.Token

	if t, err = jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			return secret, nil
		}
	}); err == nil {

		if c, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			claims = c
		} else {
			err = fmt.Errorf("token expired")
		}
	}

	return
}

func GetPublicKey(access string) (string, error) {
	if c, err := ParseAccess(access); err == nil {
		return c["public_key"].(string), nil
	}

	return "", fmt.Errorf("token invalid")
}
