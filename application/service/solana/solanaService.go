package solana

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/application/exception"
	"sharks/ports/outbound"
)

type SimpleSolanaService struct {
	solanaClient    outbound.SolanaClient
	tokenRepository outbound.TokenRepository
}

func NewSolanaService(c outbound.SolanaClient, tr outbound.TokenRepository) *SimpleSolanaService {
	return &SimpleSolanaService{c, tr}
}

func (s *SimpleSolanaService) GetNftByMint(pk string) (res *application.TokenMetadata, e *exception.Exception) {
	var err error

	if res, err = s.solanaClient.GetTokenMetadata(solana.MustPublicKeyFromBase58(pk)); err != nil {
		e = exception.FromString("Something went wrong", exception.BadRequest, exception.Unknown)
	}

	return
}

func (s *SimpleSolanaService) GetAllNftMintByOwner(pk string) (mints []string, e *exception.Exception) {
	var accounts []*token.Account
	var err error

	if accounts, err = s.solanaClient.GetTokenAccountsByWalletOwner(solana.MustPublicKeyFromBase58(pk)); err != nil {
		e = exception.FromString("Something went wrong", exception.BadRequest, exception.Unknown)

		return
	}

	mch := make(chan string)
	ech := make(chan *exception.Exception)

	for _, acc := range accounts {
		copyAcc := acc

		go func(a *token.Account) {
			if a.Amount != 0 {
				if meta := s.tokenRepository.FindByPublicKey(a.Mint.String()); meta != nil {
					if meta.IsNft {
						mch <- meta.PublicKey

						return
					}
				} else if isNft, err := s.isNft(*a); err != nil {
					ech <- exception.FromString("Something went wrong", exception.BadRequest, exception.Unknown)

					return
				} else if isNft {
					mch <- a.Mint.String()

					return
				} else {
					ft := application.TokenMetadata{PublicKey: a.Mint.String(), IsNft: false}

					if err := s.tokenRepository.Save(&ft); err != nil {
						logger.Log.Error(err.Error())
					}
				}
			}

			mch <- ""
		}(copyAcc)
	}

	for i := 0; i < len(accounts); i++ {
		select {
		case mint := <-mch:
			if mint != "" {
				mints = append(mints, mint)
			}
		case e = <-ech:
			logger.Log.Error(e.Error())
		}
	}

	return
}

func (s *SimpleSolanaService) VerifySignature(sigBase58 string, pkBase58 string, msg string) (e *exception.Exception) {
	if pk, err := solana.PublicKeyFromBase58(pkBase58); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("the owner pk is invalid", exception.Unauthorized, exception.InvalidSignature)
	} else if sig, err := solana.SignatureFromBase58(sigBase58); err != nil {
		logger.Log.Error(err.Error())

		e = exception.FromString("the signature is invalid", exception.Unauthorized, exception.InvalidSignature)
	} else if !sig.Verify(pk, []byte(msg)) {
		e = exception.FromString("the signature is invalid", exception.Unauthorized, exception.InvalidSignature)
	}

	return
}

func (s *SimpleSolanaService) isNft(t token.Account) (res bool, err error) {
	var tokenAmount *rpc.UiTokenAmount

	if tokenAmount, err = s.solanaClient.GetTokenSupply(t.Mint); err != nil {
		return
	}

	res = tokenAmount.Amount == "1" && tokenAmount.Decimals == 0

	return
}
