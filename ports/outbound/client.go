package outbound

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"sharks/application"
)

type SolanaClient interface {
	GetTokenMetadata(mint solana.PublicKey) (*application.TokenMetadata, error)
	GetTokenAccountsByWalletOwner(pk solana.PublicKey) ([]*token.Account, error)
	GetTokenSupply(mint solana.PublicKey) (*rpc.UiTokenAmount, error)
}
