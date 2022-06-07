package repository

import (
	"context"
	mng "sharks/adapters/outbound/repository/mongo"
	"sharks/application"
)

func ctx() context.Context {
	return context.Background()
}

func mapEntityToTokenMetadata(entity mng.MetadataDocument) *application.TokenMetadata {
	return &application.TokenMetadata{
		PublicKey: entity.PublicKey,
		ImageUrl:  entity.ImageUrl,
		Creators:  entity.Creators,
		IsNft:     entity.IsNft,
	}
}

func mapTokenMetadataToEntity(domain *application.TokenMetadata) *mng.MetadataDocument {
	return &mng.MetadataDocument{
		PublicKey: domain.PublicKey,
		ImageUrl:  domain.ImageUrl,
		Creators:  domain.Creators,
		IsNft:     domain.IsNft,
	}
}
