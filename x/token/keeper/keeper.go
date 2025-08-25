// In file: x/token/keeper/keeper.go

package keeper

import (
	"context" // New: Needed for collections interactions
	"fmt"
	"strings" // New: Needed for symbol normalization

	"omnis/x/token/types"

	"cosmossdk.io/collections" // Keep collections
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store" // New: Add store.KVStoreService
	"cosmossdk.io/log" // Use cosmossdk.io/log instead of tendermint/tendermint/libs/log
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Keeper struct {
	// Keeper dependencies
	cdc          codec.BinaryCodec
	addressCodec address.Codec
	storeService store.KVStoreService
	logger       log.Logger

	// Here we declare our module dependencies
	bankKeeper    banktypes.BankKeeper // Use the interface from banktypes
	accountKeeper authtypes.AccountKeeper // Use the interface from authtypes

	// Collections
	TokenSeq      collections.Sequence
	Token         collections.Map[uint64, types.Token]
	TokenBySymbol collections.Map[string, uint64] // New: Index for fast symbol lookup
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	bankKeeper banktypes.BankKeeper,
	accountKeeper authtypes.AccountKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		storeService: storeService,
		logger:       logger,
		bankKeeper:   bankKeeper,
		accountKeeper: accountKeeper,
		TokenSeq: collections.NewSequence(sb, types.TokenKeyPrefix+"_seq"),
		Token: collections.NewMap(sb, types.TokenKeyPrefix, collections.Uint64Key, codec.CollValue[types.Token](cdc)),
		TokenBySymbol: collections.NewMap(sb, types.TokenSymbolKeyPrefix, collections.StringKey, codec.CollValue[uint64](cdc)),
	}

	// This is not needed with the new collections.NewSchemaBuilder
	// schema, err := sb.Build()
	// if err != nil {
	// 	panic(err)
	// }
	// k.Schema = schema

	return k
}

// GetTokenBySymbol retrieves a token by its symbol from the store.
// This is now efficient due to the dedicated index.
func (k Keeper) GetTokenBySymbol(ctx context.Context, symbol string) (val types.Token, found bool) {
	// Normalize the symbol for case-insensitive lookup if needed
	symbol = strings.ToLower(symbol)

	// Get the token ID from the symbol-to-ID index
	id, err := k.TokenBySymbol.Get(ctx, symbol)
	if err != nil {
		return types.Token{}, false
	}

	// Get the token itself using the ID
	token, err := k.Token.Get(ctx, id)
	if err != nil {
		return types.Token{}, false
	}

	return token, true
}

// AppendToken adds a new token to the store and returns its ID.
// This function needs to be updated to also write to the new TokenBySymbol index.
func (k Keeper) AppendToken(ctx context.Context, token types.Token) uint64 {
	id, err := k.TokenSeq.Next(ctx)
	if err != nil {
		panic("failed to get next token ID")
	}
	token.Id = id

	// Set the token
	err = k.Token.Set(ctx, id, token)
	if err != nil {
		panic("failed to set token")
	}

	// Set the symbol-to-ID mapping
	err = k.TokenBySymbol.Set(ctx, strings.ToLower(token.Symbol), id)
	if err != nil {
		panic("failed to set token symbol index")
	}

	return id
}

// Logger returns the module's logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}