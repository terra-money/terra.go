package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/net/context/ctxhttp"

	"github.com/terra-money/terra.go/msg"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"

	feeutils "github.com/terra-money/core/custom/auth/client/utils"
)

// EstimateFeeResWrapper - wrapper for estimate fee query
type EstimateFeeResWrapper struct {
	Height msg.Int                  `json:"height"`
	Result feeutils.EstimateFeeResp `json:"result"`
}

// EstimateFee simulates gas and fee for a transaction
func (lcd LCDClient) EstimateFee(ctx context.Context, options CreateTxOptions) (res *feeutils.EstimateFeeResp, err error) {

	estimateReq := feeutils.EstimateFeeReq{
		BaseReq: rest.BaseReq{
			From:          msg.AccAddress(lcd.PrivKey.PubKey().Address()).String(),
			Memo:          options.Memo,
			ChainID:       lcd.ChainID,
			AccountNumber: options.AccountNumber,
			Sequence:      options.Sequence,
			TimeoutHeight: options.TimeoutHeight,
			Fees:          options.FeeAmount,
			GasPrices:     msg.NewDecCoins(lcd.GasPrice),
			Gas:           "auto",
			GasAdjustment: lcd.GasAdjustment.String(),
		},
		Msgs: options.Msgs,
	}

	reqBytes, err := lcd.EncodingConfig.Amino.MarshalJSON(estimateReq)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal")
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/txs/estimate_fee", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to estimate")
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var response EstimateFeeResWrapper
	err = lcd.EncodingConfig.Amino.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return &response.Result, nil
}
