package client

import (
	"context"
	"net/http"
	"time"

	"github.com/terra-project/terra.go/key"
	"github.com/terra-project/terra.go/msg"
	"github.com/terra-project/terra.go/tx"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	terraapp "github.com/terra-money/core/app"
	terraappparams "github.com/terra-money/core/app/params"
)

// LCDClient outer interface for building & signing & broadcasting tx
type LCDClient struct {
	URL           string
	ChainID       string
	GasPrice      msg.DecCoin
	GasAdjustment msg.Dec

	PrivKey        key.PrivKey
	EncodingConfig terraappparams.EncodingConfig

	c *http.Client
}

// NewLCDClient create new LCDClient
func NewLCDClient(URL, chainID string, gasPrice msg.DecCoin, gasAdjustment msg.Dec, tmKey key.PrivKey, httpTimeout time.Duration) *LCDClient {
	return &LCDClient{
		URL:            URL,
		ChainID:        chainID,
		GasPrice:       gasPrice,
		GasAdjustment:  gasAdjustment,
		PrivKey:        tmKey,
		EncodingConfig: terraapp.MakeEncodingConfig(),
		c:              &http.Client{Timeout: httpTimeout},
	}
}

// CreateTxOptions tx creation options
type CreateTxOptions struct {
	Msgs []msg.Msg
	Memo string

	// Optional parameters
	AccountNumber uint64
	Sequence      uint64
	GasLimit      uint64
	FeeAmount     msg.Coins

	SignMode      tx.SignMode
	FeeGranter    msg.AccAddress
	TimeoutHeight uint64
}

// CreateAndSignTx build and sign tx
func (lcd *LCDClient) CreateAndSignTx(ctx context.Context, options CreateTxOptions) (*tx.Builder, error) {
	txbuilder := tx.NewTxBuilder(lcd.EncodingConfig.TxConfig)
	txbuilder.SetFeeAmount(options.FeeAmount)
	txbuilder.SetFeeGranter(options.FeeGranter)
	txbuilder.SetGasLimit(options.GasLimit)
	txbuilder.SetMemo(options.Memo)
	txbuilder.SetMsgs(options.Msgs...)
	txbuilder.SetTimeoutHeight(options.TimeoutHeight)

	if options.FeeAmount.IsZero() || options.GasLimit == 0 {
		res, err := lcd.EstimateFee(ctx, options)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to estimate fee")
		}

		txbuilder.SetFeeAmount(res.Fee.Amount)
		txbuilder.SetGasLimit(res.Fee.Gas)
	}

	if options.AccountNumber == 0 || options.Sequence == 0 {
		account, err := lcd.LoadAccount(ctx, msg.AccAddress(lcd.PrivKey.PubKey().Address()))
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to load account")
		}

		options.AccountNumber = account.GetAccountNumber()
		options.Sequence = account.GetSequence()
	}

	// use direct sign mode as default
	if tx.SignModeUnspecified == options.SignMode {
		options.SignMode = tx.SignModeDirect
	}

	err := txbuilder.Sign(options.SignMode, tx.SignerData{
		AccountNumber: options.AccountNumber,
		ChainID:       lcd.ChainID,
		Sequence:      options.Sequence,
	}, lcd.PrivKey, true)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to sign tx")
	}

	return &txbuilder, nil
}
