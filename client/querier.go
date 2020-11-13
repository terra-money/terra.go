package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/net/context/ctxhttp"

	"github.com/terra-project/terra.go/msg"
	"github.com/terra-project/terra.go/tx"
)

// EstimateFeeReq request
type EstimateFeeReq struct {
	Tx            tx.StdTxData `json:"tx"`
	GasAdjustment string       `json:"gas_adjustment"`
	GasPrices     msg.DecCoins `json:"gas_prices"`
}

// EstimateFeeResp response
type EstimateFeeResp struct {
	Fees msg.Coins `json:"fees"`
	Gas  msg.Int   `json:"gas"`
}

// EstimateFeeResWrapper - wrapper for estimate fee query
type EstimateFeeResWrapper struct {
	Height msg.Int         `json:"height"`
	Result EstimateFeeResp `json:"result"`
}

// EstimateFee simulates gas and fee for a transaction
func (lcd LCDClient) EstimateFee(ctx context.Context, stdTx tx.StdTx) (res *EstimateFeeResp, err error) {
	broadcastReq := EstimateFeeReq{
		Tx:            stdTx.Value,
		GasAdjustment: lcd.GasAdjustment.String(),
		GasPrices:     msg.DecCoins{lcd.GasPrice},
	}

	reqBytes, err := json.Marshal(broadcastReq)
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
		return nil, fmt.Errorf("non 200 respose code %d, error: %s", resp.StatusCode, string(out))
	}

	var response EstimateFeeResWrapper
	err = json.Unmarshal(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return &response.Result, nil
}

// QueryAccountResData response
type QueryAccountResData struct {
	Address       msg.AccAddress `json:"address"`
	Coins         msg.Coins      `json:"coins"`
	AccountNumber msg.Int        `json:"account_number"`
	Sequence      msg.Int        `json:"sequence"`
}

// QueryAccountRes response
type QueryAccountRes struct {
	Type  string              `json:"type"`
	Value QueryAccountResData `json:"value"`
}

// QueryAccountResWrapper - wrapper for estimate fee query
type QueryAccountResWrapper struct {
	Height msg.Int         `json:"height"`
	Result QueryAccountRes `json:"result"`
}

// LoadAccount simulates gas and fee for a transaction
func (lcd LCDClient) LoadAccount(ctx context.Context, address msg.AccAddress) (res *QueryAccountResData, err error) {
	resp, err := ctxhttp.Get(ctx, lcd.c, lcd.URL+fmt.Sprintf("/auth/accounts/%s", address))
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

	var response QueryAccountResWrapper
	err = json.Unmarshal(out, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return &response.Result.Value, nil
}
