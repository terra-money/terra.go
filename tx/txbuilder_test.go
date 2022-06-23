package tx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/terra.go/key"
	"github.com/terra-money/terra.go/msg"

	terraapp "github.com/terra-money/core/v2/app"
)

func Test_Sign(t *testing.T) {
	mnemonic := "essence gallery exit illegal nasty luxury sport trouble measure benefit busy almost bulb fat shed today produce glide meadow require impact fruit omit weasel"
	privKeyBz, err := key.DerivePrivKeyBz(mnemonic, key.CreateHDPath(0, 0))
	assert.NoError(t, err)
	privKey, err := key.PrivKeyGen(privKeyBz)
	assert.NoError(t, err)

	addr := msg.AccAddress(privKey.PubKey().Address())
	assert.Equal(t, addr.String(), "terra1cevwjzwft3pjuf5nc32d9kyrvh5y7fp9havw7k")

	txBuilder := NewTxBuilder(terraapp.MakeEncodingConfig().TxConfig)
	err = txBuilder.SetMsgs(
		msg.NewMsgExecuteContract(
			addr,
			addr,
			[]byte("{\"withdraw\":{\"position_idx\":\"1\",\"collateral\":{\"info\":{\"native_token\":{\"denom\":\"uusd\"}},\"amount\":\"1000\"}}}"),
			msg.Coins{},
		),
	)
	require.NoError(t, err)
	txBuilder.SetFeeAmount(msg.Coins{})
	txBuilder.SetGasLimit(1000000)

	// amino version test
	err = txBuilder.Sign(SignModeLegacyAminoJSON, SignerData{
		ChainID:       "testnet",
		AccountNumber: 359,
		Sequence:      4,
	}, privKey, true)
	require.NoError(t, err)

	sigs, err := txBuilder.GetTx().GetSignaturesV2()
	require.NoError(t, err)

	bz, err := txBuilder.TxConfig.MarshalSignatureJSON(sigs)
	fmt.Println(string(bz))
	assert.NoError(t, err)
	assert.Equal(t, bz, []byte(`{"signatures":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AmADjpxwusAnJ7ahD7+trzovH32w+LaRGVZSZUOd3E3d"},"data":{"single":{"mode":"SIGN_MODE_LEGACY_AMINO_JSON","signature":"FUoV2W4aS8zm2AzmrduvZCw8QXuAZXz/hgp/aCr/jtwWi6oHpsMhhR+dUt1r0L29PAUJz69aMVvMfTQEecI0+w=="}},"sequence":"4"}]}`))

	// direct mode test
	err = txBuilder.Sign(SignModeDirect, SignerData{
		ChainID:       "testnet",
		AccountNumber: 359,
		Sequence:      4,
	}, privKey, true)
	require.NoError(t, err)

	sigs, err = txBuilder.GetTx().GetSignaturesV2()
	require.NoError(t, err)

	bz, err = txBuilder.TxConfig.MarshalSignatureJSON(sigs)
	assert.NoError(t, err)
	assert.Equal(t, bz, []byte(`{"signatures":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AmADjpxwusAnJ7ahD7+trzovH32w+LaRGVZSZUOd3E3d"},"data":{"single":{"mode":"SIGN_MODE_DIRECT","signature":"BllKHpAIgMHVvip11+QAHqlhGGfh6dIpFAdIqE+EZFIs4oWggjjua9g9bSZY9Sr6y+8Fjn2X7ziWUHT8zSbNdQ=="}},"sequence":"4"}]}`))
}
