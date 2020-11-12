package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-project/terra.go/msg"
	"github.com/terra-project/terra.go/tx"
)

// EstimateFeeReq request
type EstimateFeeReq struct {
	Tx            tx.StdTxData `json:"tx"`
	GasAdjustment string       `json:"gas_adjustment"`
	GasPrices     sdk.DecCoins `json:"gas_prices"`
}

// EstimateFeeResp response
type EstimateFeeResp struct {
	Fees msg.Coins `json:"fees"`
	Gas  msg.Int   `json:"gas"`
}

// EstimateFeeResponseWrapper - wrapper for estimate fee query
type EstimateFeeResponseWrapper struct {
	Height msg.Int         `json:"height"`
	Result EstimateFeeResp `json:"result"`
}

// EstimateFee simulates gas and fee for a transaction
func (lcdClient LCDClient) EstimateFee(stdTx tx.StdTx) (res EstimateFeeResp, err error) {
	broadcastReq := EstimateFeeReq{
		Tx:            stdTx.Value,
		GasAdjustment: lcdClient.GasAdjustment.String(),
		GasPrices:     msg.DecCoins{lcdClient.GasPrice},
	}

	reqBytes, err := json.Marshal(broadcastReq)
	if err != nil {
		return EstimateFeeResp{}, sdkerrors.Wrap(err, "failed to marshal")
	}

	resp, err := http.Post(lcdClient.URL+"/txs/estimate_fee", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return EstimateFeeResp{}, sdkerrors.Wrap(err, "failed to estimate")
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EstimateFeeResp{}, sdkerrors.Wrap(err, "failed to read response")
	}

	var response EstimateFeeResponseWrapper
	err = json.Unmarshal(out, &response)
	if err != nil {
		return EstimateFeeResp{}, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return response.Result, nil
}
