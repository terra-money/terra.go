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

// BroadcastReq broadcast request body
type BroadcastReq struct {
	Tx   tx.StdTxData `json:"tx"`
	Mode string       `json:"mode"`
}

// TxResponse response
type TxResponse struct {
	Height msg.Int `json:"height"`
	TxHash string  `json:"txhash"`
	Code   uint32  `json:"code,omitempty"`
	RawLog string  `json:"raw_log,omitempty"`
}

// Broadcast - no-lint
func (lcd LCDClient) Broadcast(ctx context.Context, stdTx *tx.StdTx) (*TxResponse, error) {
	broadcastReq := BroadcastReq{
		Tx:   stdTx.Value,
		Mode: "sync",
	}

	reqBytes, err := json.Marshal(broadcastReq)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal")
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/txs", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to broadcast")
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response code %d: %s", resp.StatusCode, string(out))
	}

	var txResponse TxResponse
	err = json.Unmarshal(out, &txResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	if txResponse.Code != 0 {
		return &txResponse, fmt.Errorf("tx failed with code %d: %s", txResponse.Code, txResponse.RawLog)
	}

	return &txResponse, nil
}
