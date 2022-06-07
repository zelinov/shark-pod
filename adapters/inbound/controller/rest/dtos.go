package rest

import "time"

type TokenCreated struct {
	Access    string    `json:"accessToken"`
	Refresh   string    `json:"refreshToken"`
	PublicKey string    `json:"publicKey"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type LoginRequest struct {
	Signature string `json:"nonce"`
	PublicKey string `json:"publicKey"`
}

type LoginResponse TokenCreated

type RefreshRequest struct {
	Refresh string `json:"refresh"`
}

type RefreshResponse TokenCreated

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type CheckResponse struct {
	PublicKey string `json:"publicKey"`
}

type GetAllTokensByOwnerResponse struct {
	Images []ImageDto `json:"images"`
}

type ImageDto struct {
	Url string `json:"url"`
}
