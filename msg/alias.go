package msg

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type (
	// Msg sdk.Msg
	Msg = sdk.Msg

	// Send bank send msg
	Send = banktypes.MsgSend

	// MultiSend bank multi-send msg
	MultiSend = banktypes.MsgMultiSend

	// StoreCode wasm store_code msg
	StoreCode = wasmtypes.MsgStoreCode

	// InstantiateContract wasm contract initiation msg
	InstantiateContract = wasmtypes.MsgInstantiateContract

	// ExecuteContract wasm contract execution msg
	ExecuteContract = wasmtypes.MsgExecuteContract

	// MigrateContract wasm contract migration msg
	MigrateContract = wasmtypes.MsgMigrateContract

	// AccessConfig wasm contract access config
	AccessConfig = wasmtypes.AccessConfig

	// Coin nolint
	Coin = sdk.Coin
	// Coins nolint
	Coins = sdk.Coins
	// DecCoin nolint
	DecCoin = sdk.DecCoin
	// DecCoins nolint
	DecCoins = sdk.DecCoins

	// Int nolint
	Int = sdk.Int
	// Dec nolint
	Dec = sdk.Dec

	// AccAddress nolint
	AccAddress = sdk.AccAddress
	// ValAddress nolint
	ValAddress = sdk.ValAddress
	// ConsAddress nolint
	ConsAddress = sdk.ConsAddress
)

// function alias
var (
	NewMsgSend                = banktypes.NewMsgSend
	NewMsgMultiSend           = banktypes.NewMsgMultiSend

	NewCoin         = sdk.NewCoin
	NewInt64Coin    = sdk.NewInt64Coin
	NewCoins        = sdk.NewCoins
	NewDecCoin      = sdk.NewDecCoin
	NewInt64DecCoin = sdk.NewInt64DecCoin
	NewDecCoins     = sdk.NewDecCoins

	NewInt                   = sdk.NewInt
	NewIntFromBigInt         = sdk.NewIntFromBigInt
	NewIntFromString         = sdk.NewIntFromString
	NewIntFromUint64         = sdk.NewIntFromUint64
	NewIntWithDecimal        = sdk.NewIntWithDecimal
	NewDec                   = sdk.NewDec
	NewDecCoinFromCoin       = sdk.NewDecCoinFromCoin
	NewDecCoinFromDec        = sdk.NewDecCoinFromDec
	NewDecFromBigInt         = sdk.NewDecFromBigInt
	NewDecFromBigIntWithPrec = sdk.NewDecFromBigIntWithPrec
	NewDecFromInt            = sdk.NewDecFromInt
	NewDecFromIntWithPrec    = sdk.NewDecFromIntWithPrec
	NewDecFromStr            = sdk.NewDecFromStr
	NewDecWithPrec           = sdk.NewDecWithPrec
	AccAddressFromBech32     = sdk.AccAddressFromBech32
	AccAddressFromHex        = sdk.AccAddressFromHex
	ValAddressFromBech32     = sdk.ValAddressFromBech32
	ValAddressFromHex        = sdk.ValAddressFromHex
	ConsAddressFromBech32    = sdk.ConsAddressFromBech32
	ConsAddressFromHex       = sdk.ConsAddressFromHex
)
