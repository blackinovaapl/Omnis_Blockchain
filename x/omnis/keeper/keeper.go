package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"omnis/x/omnis/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]
	 Items collections.Map[string, types.Item] 

// Inside type Keeper struct { ... }

// SetItem stores an item in the KVStore.
func (k Keeper) SetItem(ctx context.Context, item types.Item) error {
    return k.Items.Set(ctx, item.Id, item) // Assuming item.Id is your string key
}

// GetItem retrieves an item from the KVStore.
func (k Keeper) GetItem(ctx context.Context, id string) (types.Item, error) {
    item, err := k.Items.Get(ctx, id)
    if err != nil {
        return types.Item{}, err // Handle collections.ErrNotFound in calling code
    }
    return item, nil
}

// HasItem checks if an item exists in the KVStore.
func (k Keeper) HasItem(ctx context.Context, id string) (bool, error) {
    return k.Items.Has(ctx, id)
}

// DeleteItem deletes an item from the KVStore.
func (k Keeper) DeleteItem(ctx context.Context, id string) error {
    return k.Items.Remove(ctx, id)
}

// IterateItems iterates over all items.
func (k Keeper) IterateItems(ctx context.Context, cb func(id string, item types.Item) (bool, error)) error {
    return k.Items.Walk(ctx, nil, cb)
}

}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	 Items: collections.NewMap(
            sb,
            types.ItemKeyPrefix, // Unique prefix for this collection
            "items",             // Human-readable name for the collection
            collections.StringKey, // Key type (e.g., string, uint64, sdk.AccAddress)
            codec.CollValue[types.Item](cdc), // Value type (your protobuf message)
        ),
	
	
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
