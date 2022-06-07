package client

import "github.com/gagliardetto/solana-go"

type MetaplexMeta struct {
	Key             int8
	UpdateAuthority solana.PublicKey
	Mint            solana.PublicKey
	Data            MetaplexData
}

type MetaplexData struct {
	Name   string
	Symbol string
	Uri    string
}
