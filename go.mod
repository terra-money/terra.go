module github.com/terra-money/terra.go

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/go-bip39 v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/terra-money/core v0.5.13
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76

