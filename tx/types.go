package tx

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type (

	// Builder to create transaction for broadcasting
	Builder struct {
		client.TxBuilder
		client.TxConfig
	}

	// SignerData is the specific information needed to sign a transaction that generally
	// isn't included in the transaction body itself
	SignerData = authsigning.SignerData

	// SignMode represents a signing mode with its own security guarantees.
	SignMode = signing.SignMode

	// SignatureV2 is a convenience type that is easier to use in application logic
	// than the protobuf SignerInfo's and raw signature bytes. It goes beyond the
	// first sdk.Signature types by supporting sign modes and explicitly nested
	// multi-signatures. It is intended to be used for both building and verifying
	// signatures.
	SignatureV2 = signing.SignatureV2

	// Tx defines a transaction interface that supports all standard message, signature
	// fee, memo, and auxiliary interfaces.
	Tx = authsigning.Tx
)

const (
	// SignModeUnspecified specifies an unknown signing mode and will be
	// rejected
	SignModeUnspecified SignMode = 0
	// SignModeDirect specifies a signing mode which uses SignDoc and is
	// verified with raw bytes from Tx
	SignModeDirect SignMode = 1
	// SignModeTexture is a future signing mode that will verify some
	// human-readable textual representation on top of the binary representation
	// from SIGN_MODE_DIRECT
	SignModeTexture SignMode = 2
	// SignModeLegacyAminoJSON is a backwards compatibility mode which uses
	// Amino JSON and will be removed in the future
	SignModeLegacyAminoJSON SignMode = 127
)
