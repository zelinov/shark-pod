package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"net/http"
	"sharks/adapters/outbound/logger"
	"sharks/application"
	"sharks/config"
	"strings"
	"time"
)

type SimpleSolanaClient struct {
	client rpc.Client
}

func NewSimpleSolanaClient() *SimpleSolanaClient {
	return &SimpleSolanaClient{
		client: *rpc.New(config.GetConfig().SolanaRpcPoolUrl),
	}
}

var splTokenProgramPublicKey = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
var metadataAccountPublicKey = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
var candyMachineV2ProgramID = solana.MustPublicKeyFromBase58("cndyAnrLdpjq1Ssp1z8xxDsB8dxe7u4HL5Nxi2K5WXZ")
var metadataNameByteArray = []byte("metadata")
var metadataAccountPublicKeyByteArray = metadataAccountPublicKey.Bytes()

func (c *SimpleSolanaClient) GetTokenAccountsByWalletOwner(pk solana.PublicKey) (accounts []*token.Account, err error) {
	conf := rpc.GetTokenAccountsConfig{
		ProgramId: &splTokenProgramPublicKey,
	}
	opts := rpc.GetTokenAccountsOpts{
		Encoding: solana.EncodingBase64,
	}

	var out *rpc.GetTokenAccountsResult

	if out, err = c.client.GetTokenAccountsByOwner(context.Background(), pk, &conf, &opts); err != nil {
		logger.Log.Error(fmt.Sprintf("GetTokenAccountsByOwner for owner %s, err: %s", pk.String(), err.Error()))

		return
	}

	for _, rawAccount := range out.Value {
		var tokAcc token.Account

		if err = bin.NewBinDecoder(rawAccount.Account.Data.GetBinary()).Decode(&tokAcc); err != nil {
			logger.Log.Error(fmt.Sprintf("DecodeTokenData for owner %s, err: %s", pk.String(), err.Error()))

			return
		}

		accounts = append(accounts, &tokAcc)
	}

	return
}

func (c *SimpleSolanaClient) GetTokenSupply(mint solana.PublicKey) (tokenAmount *rpc.UiTokenAmount, err error) {
	var out *rpc.GetTokenSupplyResult

	if out, err = c.client.GetTokenSupply(context.Background(), mint, rpc.CommitmentFinalized); err == nil {
		tokenAmount = out.Value
	} else {
		logger.Log.Error(fmt.Sprintf("GetTokenSupply for mint %s, err: %s", mint.String(), err.Error()))
	}

	return
}

func (c *SimpleSolanaClient) GetTokenMetadata(mint solana.PublicKey) (tokenMeta *application.TokenMetadata, err error) {
	var metaData []byte
	var meta MetaplexMeta
	var data map[string]interface{}

	if metaData, err = c.getAccountMetaData(mint); err != nil {
		logger.Log.Error(fmt.Sprintf("getAccountMetaData for mint %s, err: %s", mint.String(), err.Error()))
	} else if meta, err = getMetaplexMeta(metaData); err != nil {
		logger.Log.Error(fmt.Sprintf("getMetaplexMeta for mint %s, err: %s", mint.String(), err.Error()))
	} else if data, err = getTokenAdditionalData(meta.Data.Uri); err != nil {
		logger.Log.Error(fmt.Sprintf("getTokenAdditionalData for mint %s, err: %s", mint.String(), err.Error()))
	} else {
		var creatorsPk []string

		if props, ok := data["properties"].(map[string]interface{}); ok {
			if cr, ok := props["creators"]; ok {
				if creators, ok := cr.([]interface{}); ok {
					for _, creator := range creators {
						creatorsPk = append(creatorsPk, creator.(map[string]interface{})["address"].(string))
					}
				}
			}
		}

		if candy, e := c.getCandyMachineCreator(mint); e == nil {
			if candy != nil {
				creatorsPk = append(creatorsPk, (*candy).String())
			}
		} else {
			logger.Log.Error(fmt.Sprintf("getCandyMachineCreator for mint %s, err: %s", mint.String(), e.Error()))
		}

		tokenMeta = &application.TokenMetadata{
			PublicKey: mint.String(),
			ImageUrl:  data["image"].(string),
			Creators:  creatorsPk,
			IsNft:     true,
		}
	}

	return
}

func (c *SimpleSolanaClient) getAccountMetaData(pk solana.PublicKey) (data []byte, err error) {
	pdaSeeds := [][]byte{
		metadataNameByteArray,
		metadataAccountPublicKeyByteArray,
		pk.Bytes(),
	}

	var pdaAddr solana.PublicKey

	if pdaAddr, _, err = solana.FindProgramAddress(pdaSeeds, metadataAccountPublicKey); err != nil {
		return
	}

	var out *rpc.GetAccountInfoResult

	if out, err = c.client.GetAccountInfo(context.Background(), pdaAddr); err == nil {
		data = out.Value.Data.GetBinary()
	}

	return
}

var client = http.Client{
	Timeout: 5000 * time.Millisecond,
}

func getTokenAdditionalData(uri string) (data map[string]interface{}, err error) {
	var resp *http.Response

	if resp, err = client.Get(uri); err == nil {
		if (resp.StatusCode / 100) != 2 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			err = fmt.Errorf("reqest %s. response [status: %d, body: %s]", uri, resp.StatusCode, buf.String())
		} else {
			err = json.NewDecoder(resp.Body).Decode(&data)
		}
	}

	return
}

func getMetaplexMeta(data []byte) (meta MetaplexMeta, err error) {
	if err = bin.NewBorshDecoder(data).Decode(&meta); err == nil {
		normalize(&meta)
	}

	return
}

func (c *SimpleSolanaClient) getCandyMachineCreator(mint solana.PublicKey) (cmId *solana.PublicKey, err error) {
	if sig, e := c.getFirstSignature(mint, nil); err != nil {
		err = e
	} else if tx, e := c.client.GetTransaction(context.Background(), *sig, nil); e != nil {
		err = e
	} else {
		msg := tx.Transaction.GetParsedTransaction().Message

		for _, ins := range msg.Instructions {
			if msg.AccountKeys[ins.ProgramIDIndex] == candyMachineV2ProgramID {
				cmId = &msg.AccountKeys[ins.Accounts[1]]
			}
		}
	}

	return
}

func (c *SimpleSolanaClient) getFirstSignature(
	mint solana.PublicKey,
	before *solana.Signature,
) (sig *solana.Signature, err error) {
	opts := rpc.GetSignaturesForAddressOpts{}

	if before != nil {
		opts.Before = *before
	}

	if sigs, e := c.client.GetSignaturesForAddressWithOpts(context.Background(), mint, &opts); err != nil {
		err = e
	} else if len(sigs) == 1000 {
		sig, err = c.getFirstSignature(mint, &sigs[len(sigs)-1].Signature)
	} else if len(sigs) == 0 {
		sig = before
	} else {
		sig = &sigs[len(sigs)-1].Signature
	}

	return
}

func normalize(meta *MetaplexMeta) {
	meta.Data.Uri = strings.ReplaceAll(meta.Data.Uri, "\u0000", "")
	meta.Data.Symbol = strings.ReplaceAll(meta.Data.Symbol, "\u0000", "")
	meta.Data.Name = strings.ReplaceAll(meta.Data.Name, "\u0000", "")
}
