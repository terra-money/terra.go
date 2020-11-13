package client

import (
	"context"
	"net/http"
	"time"

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

	c *http.Client
}

// NewLCDClient create new LCDClient
func NewLCDClient(URL, chainID string, gasPrice msg.DecCoin, gasAdjustment msg.Dec, tmKey key.StdPrivKey, httpTimeout time.Duration) *LCDClient {
	return &LCDClient{
		URL:           URL,
		ChainID:       chainID,
		GasPrice:      gasPrice,
		GasAdjustment: gasAdjustment,
		TmKey:         tmKey,
		c:             &http.Client{Timeout: httpTimeout},
	}
}

// CreateTxOptions tx creation options
type CreateTxOptions struct {
	Msgs []msg.Msg
	Memo string

	// Optional parameters
	AccountNumber msg.Int
	Sequence      msg.Int
	Fee           tx.StdFee
}

// CreateAndSignTx build and sign tx
func (lcd *LCDClient) CreateAndSignTx(ctx context.Context, options CreateTxOptions) (*tx.StdTx, error) {
	stdTx := tx.NewStdTx(options.Msgs, options.Memo, options.Fee)
	if options.Fee.IsEmpty() {
		fee, err := lcd.EstimateFee(ctx, stdTx)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to estimate fee")
		}

		stdTx.Value.Fee.Amount = fee.Fees
		stdTx.Value.Fee.Gas = fee.Gas
	}

	if (msg.Int{}) == options.AccountNumber ||
		(msg.Int{}) == options.Sequence ||
		options.AccountNumber.IsZero() {
		account, err := lcd.LoadAccount(ctx, msg.AccAddress(lcd.TmKey.PubKey().Address()))
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to load account")
		}

		options.AccountNumber = account.AccountNumber
		options.Sequence = account.Sequence
	}

	signature, err := stdTx.Sign(lcd.TmKey, lcd.ChainID, options.AccountNumber, options.Sequence)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to sign tx")
	}

	stdTx.AppendSignatures(signature)
	return &stdTx, nil
}
