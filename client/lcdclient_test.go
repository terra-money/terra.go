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
	mnemonic := "satisfy adjust timber high purchase tuition stool faith fine install that you unaware feed domain license impose boss human eager hat rent enjoy dawn"
	privKeyBz, err := key.DerivePrivKeyBz(mnemonic, key.CreateHDPath(0, 0))
	assert.NoError(t, err)
	privKey, err := key.PrivKeyGen(privKeyBz)
	assert.NoError(t, err)

	addr := msg.AccAddress(privKey.PubKey().Address())
	assert.Equal(t, addr.String(), "terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8")

	toAddr, err := msg.AccAddressFromBech32("terra1t849fxw7e8ney35mxemh4h3ayea4zf77dslwna")
	assert.NoError(t, err)

	LCDClient := NewLCDClient(
		"http://127.0.0.1:1317",
		"testnet",
		msg.NewDecCoinFromDec("uluna", msg.NewDecFromIntWithPrec(msg.NewInt(15), 2)), // 0.15uluna
		msg.NewDecFromIntWithPrec(msg.NewInt(15), 1), privKey,
		10*time.Second,
	)

	tx, err := LCDClient.CreateAndSignTx(
		context.Background(),
		CreateTxOptions{
			Msgs: []msg.Msg{
				msg.NewMsgSend(addr, toAddr, msg.NewCoins(msg.NewInt64Coin("uluna", 1000000))), // 1Luna
			},
			Memo:     "",
			SignMode: tx.SignModeDirect,
		})
	assert.NoError(t, err)

	res, err := LCDClient.Broadcast(context.Background(), tx)
	assert.NoError(t, err)
	assert.Equal(t, res.Code, uint32(0))
}
