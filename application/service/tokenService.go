package service

import (
	"golang.org/x/exp/slices"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/application/exception"
	"sharks/ports/inbound"
	"sharks/ports/outbound"
)

type NftTokenService struct {
	nftRepository outbound.TokenRepository
	solanaService inbound.SolanaService
}

func NewNftTokenService(r outbound.TokenRepository, s inbound.SolanaService) *NftTokenService {
	return &NftTokenService{r, s}
}

func (s *NftTokenService) GetAll(filter *application.TokenFilter) (tokens []*application.TokenMetadata, e *exception.Exception) {
	var tokensForSave []*application.TokenMetadata
	var mints []string

	if mints, e = s.solanaService.GetAllNftMintByOwner(filter.Owner); e == nil {
		dbtch := make(chan *application.TokenMetadata, len(mints))
		stch := make(chan *application.TokenMetadata, len(mints))
		ech := make(chan error, len(mints))

		for _, mint := range mints {
			go func(m string) {
				if t := s.nftRepository.FindByPublicKey(m); t != nil {
					dbtch <- t
				} else if t, ex := s.solanaService.GetNftByMint(m); ex == nil && t != nil {
					stch <- t
				} else if ex != nil {
					ech <- ex
				} else {
					stch <- nil
				}
			}(mint)
		}

		for i := 0; i < len(mints); i++ {
			select {
			case t := <-dbtch:
				if filter.Creators == nil || isCreatorNft(t, *filter.Creators) {
					tokens = append(tokens, t)
				}
			case t := <-stch:
				if t != nil {
					if filter.Creators == nil || isCreatorNft(t, *filter.Creators) {
						tokens = append(tokens, t)
					}

					tokensForSave = append(tokensForSave, t)
				}
			case err := <-ech:
				logger.Log.Error(err.Error())
			}
		}
	}

	if len(tokensForSave) > 0 {
		if err := s.nftRepository.SaveAll(tokensForSave); err != nil {
			logger.Log.Error(err.Error())
		}
	}

	return
}

func isCreatorNft(acc *application.TokenMetadata, creators []string) bool {
	for _, c := range acc.Creators {
		if slices.Contains(creators, c) {
			return true
		}
	}

	return false
}
