package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/terra-money/terra.go/key"
	"github.com/terra-money/terra.go/msg"
	"github.com/terra-money/terra.go/tx"
)

func Test_Transaction(t *testing.T) {
	mnemonic := "essence gallery exit illegal nasty luxury sport trouble measure benefit busy almost bulb fat shed today produce glide meadow require impact fruit omit weasel"
	privKeyBz, err := key.DerivePrivKeyBz(mnemonic, key.CreateHDPath(0, 0))
	assert.NoError(t, err)
	privKey, err := key.PrivKeyGen(privKeyBz)
	assert.NoError(t, err)

	addr := msg.AccAddress(privKey.PubKey().Address())
	assert.Equal(t, addr.String(), "terra1cevwjzwft3pjuf5nc32d9kyrvh5y7fp9havw7k")

	toAddr, err := msg.AccAddressFromBech32("terra1t849fxw7e8ney35mxemh4h3ayea4zf77dslwna")
	assert.NoError(t, err)

	LCDClient := NewLCDClient(
		"http://127.0.0.1:1317",
		"testnet",
		msg.NewDecCoinFromDec("uusd", msg.NewDecFromIntWithPrec(msg.NewInt(15), 2)), // 0.15uusd
		msg.NewDecFromIntWithPrec(msg.NewInt(15), 1), privKey,
		10*time.Second,
	)

	tx, err := LCDClient.CreateAndSignTx(
		context.Background(),
		CreateTxOptions{
			Msgs: []msg.Msg{
				msg.NewMsgSend(addr, toAddr, msg.NewCoins(msg.NewInt64Coin("uusd", 100000000))), // 100UST
				msg.NewMsgSwapSend(addr, toAddr, msg.NewInt64Coin("uusd", 1000000), "ukrw"),
			},
			Memo:     "",
			SignMode: tx.SignModeDirect,
		})
	assert.NoError(t, err)

	res, err := LCDClient.Broadcast(context.Background(), tx)
	assert.NoError(t, err)
	assert.Equal(t, res.Code, uint32(0))
}
