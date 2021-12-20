package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/net/context/ctxhttp"

	"github.com/terra-money/terra.go/tx"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// Broadcast - no-lint
func (lcd LCDClient) Broadcast(ctx context.Context, txbuilder *tx.Builder) (*sdk.TxResponse, error) {
	txBytes, err := txbuilder.GetTxBytes()
	if err != nil {
		return nil, err
	}

	broadcastReq := txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
	}

	reqBytes, err := json.Marshal(broadcastReq)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal")
	}

	resp, err := ctxhttp.Post(ctx, lcd.c, lcd.URL+"/cosmos/tx/v1beta1/txs", "application/json", bytes.NewBuffer(reqBytes))
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

	var broadcastTxResponse txtypes.BroadcastTxResponse
	err = lcd.EncodingConfig.Marshaler.UnmarshalJSON(out, &broadcastTxResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	txResponse := broadcastTxResponse.TxResponse
	if txResponse.Code != 0 {
		return txResponse, fmt.Errorf("tx failed with code %d: %s", txResponse.Code, txResponse.RawLog)
	}

	return txResponse, nil
}
