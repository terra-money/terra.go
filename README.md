# terra.go
Simple transaction build & signing library

## How to use
```
// Create Mnemonic
mnemonic, err := key.CreateMnemonic()
assert.NoError(t, err)

// Derive Raw Private Key
privKeyBz, err := key.DerivePrivKeyBz(mnemonic, key.CreateHDPath(0, 0))
assert.NoError(t, err)

// Generate StdPrivKey
privKey, err := key.PrivKeyGen(privKey)
assert.NoError(t, err)

// Generate Address from Public Key
addr := msg.AccAddress(privKey.PubKey().Address())
assert.Equal(t, addr.String(), "terra1cevwjzwft3pjuf5nc32d9kyrvh5y7fp9havw7k")

// Create LCDClient
LCDClient := NewLCDClient(
    "http://127.0.0.1:1317",
    "testnet",
    msg.NewDecCoinFromDec("uusd", msg.NewDecFromIntWithPrec(msg.NewInt(15), 2)), // 0.15uusd
    msg.NewDecFromIntWithPrec(msg.NewInt(15), 1), privKey,
)

// Create tx
tx, err := LCDClient.CreateAndSignTx(CreateTxOptions{
    Msgs: []msg.Msg{
        msg.NewSend(addr, toAddr, msg.NewCoins(msg.NewInt64Coin("uusd", 100000000))), // 100UST
    },
    Memo:          "",

    // Options Paramters (if empty, load chain info)
    // AccountNumber: msg.NewInt(33),
    // Sequence:      msg.NewInt(1),
    // Options Paramters (if empty, simulate gas & fee)
    // FeeAmount: msg.NewCoins(),
    // GasLimit: 1000000,
    // FeeGranter: msg.AccAddress{},
    // SignMode: tx.SignModeDirect, 
})
assert.NoError(t, err)

// Broadcast
res, err := LCDClient.Broadcast(context.Background(), tx)
assert.NoError(t, err)
fmt.Println(res)
```
