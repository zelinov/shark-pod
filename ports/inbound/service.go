package inbound

import (
	"sharks/application"
	"sharks/application/exception"
)

type AuthService interface {
	Nonce(pk string) *application.Nonce
	Login(cred application.Credentials) (*application.JwtToken, *exception.Exception)
	Verify(token string) *exception.Exception
	Refresh(token string) (*application.JwtToken, *exception.Exception)
	Logout(token string) *exception.Exception
}

type SolanaService interface {
	VerifySignature(sigBase58 string, pkBase58 string, msg string) *exception.Exception
	GetNftByMint(pk string) (*application.TokenMetadata, *exception.Exception)
	GetAllNftMintByOwner(pk string) ([]string, *exception.Exception)
}

type TokenService interface {
	GetAll(filter *application.TokenFilter) ([]*application.TokenMetadata, *exception.Exception)
}
