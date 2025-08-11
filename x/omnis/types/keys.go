package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "omnis"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"

    // Event types
    EventTypeCreateItem = "create_item" // New event type for item creation
    EventTypeUpdateItem = "update_item" // New event type for item update
    EventTypeDeleteItem = "delete_item" // New event type for item deletion

    // Attribute keys for events
    AttributeKeyItemID      = "item_id"
    AttributeKeyItemName    = "item_name"
    AttributeKeyNewItemName = "new_item_name" // If you allow name updates
    AttributeKeyCreator     = "creator"
    // Add other attribute keys as needed for your module's events
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_omnis")

// ItemKeyPrefix is the prefix to retrieve all Items
var ItemKeyPrefix = collections.NewPrefix("i_omnis_item") // NEW: Unique prefix for your items
// You can use bytes (e.g., collections.NewPrefix(0x01)) or strings.
// String prefixes are often more human-readable in debug tools.
// Ensure "i_omnis_item" is truly unique within your module's store.
// If you have multiple distinct data types, give each its own prefix.