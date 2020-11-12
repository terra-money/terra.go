package client

import (
	"fmt"

	"github.com/terra-project/terra.go/key"
	"github.com/terra-project/terra.go/msg"
	"github.com/terra-project/terra.go/tx"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// LCDClient outer interface for building & signing & broadcasting tx
type LCDClient struct {
	URL           string
	ChainID       string
	GasPrice      msg.DecCoin
	GasAdjustment msg.Dec

	TmKey key.StdPrivKey
}

// NewLCDClient create new LCDClient
func NewLCDClient(URL, chainID string, gasPrice msg.DecCoin, gasAdjustment msg.Dec, tmKey key.StdPrivKey) LCDClient {
	return LCDClient{
		URL:           URL,
		ChainID:       chainID,
		GasPrice:      gasPrice,
		GasAdjustment: gasAdjustment,
		TmKey:         tmKey,
	}
}

// CreateTxOptions tx creation options
type CreateTxOptions struct {
	Msgs []msg.Msg
	Fee  tx.StdFee
	Memo string

	AccountNumber msg.Int
	Sequence      msg.Int
}

// CreateAndSignTx build and sign tx
func (lcdClient LCDClient) CreateAndSignTx(options CreateTxOptions) (tx.StdTx, error) {
	stdTx := tx.NewStdTx(options.Msgs, options.Memo, options.Fee)
	if options.Fee.IsEmpty() {
		fee, err := lcdClient.EstimateFee(stdTx)
		if err != nil {
			return tx.StdTx{}, sdkerrors.Wrap(err, "failed to estimate fee")
		}

		stdTx.Value.Fee.Amount = fee.Fees
		stdTx.Value.Fee.Gas = fee.Gas
	}

	fmt.Println(stdTx)

	signature, err := stdTx.Sign(lcdClient.TmKey, lcdClient.ChainID, options.AccountNumber, options.Sequence)
	if err != nil {
		return tx.StdTx{}, sdkerrors.Wrap(err, "failed to sign tx")
	}

	stdTx.AppendSignatures(signature)
	return stdTx, nil
}
