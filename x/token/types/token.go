
package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Token represents a custom OMS-20 token on the Omnis blockchain.
type Token struct {
	// Denom is the unique identifier for the token (e.g., "usd").
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`

	// Creator is the address of the account that created the token.
	Creator string `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`

	// Supply is the total supply of the token.
	Supply sdk.Int `protobuf:"bytes,3,opt,name=supply,proto3" json:"supply,omitempty"`

	// Metadata is a map for storing custom, Solana-like metadata.
	// This can include fields like Name, Symbol, Decimals, Description, etc.
	Metadata map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty"`
}

// NewToken creates a new Token object
func NewToken(denom string, creator sdk.AccAddress, supply sdk.Int, metadata map[string]string) Token {
	return Token{
		Denom:    denom,
		Creator:  creator.String(),
		Supply:   supply,
		Metadata: metadata,
	}
}

// Validate validates the basic properties of a Token.
func (t Token) Validate() error {
	if strings.TrimSpace(t.Denom) == "" {
		return fmt.Errorf("token denom cannot be empty")
	}
	if !sdk.IsValidCoins(sdk.NewCoins(sdk.NewCoin(t.Denom, t.Supply))) {
		return fmt.Errorf("invalid token supply for denom %s: %s", t.Denom, t.Supply.String())
	}
	if _, err := sdk.AccAddressFromBech32(t.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// Add more validation for metadata if specific fields are required
	return nil
}

// ------------------------------------------------------------
// MsgCreateToken (Example Message for creating a token)
// This is a placeholder; you will typically define your messages
// in a separate `tx.go` or `msgs.go` file within the `types` package.
// But for a basic structure, it's good to see how it would look.
// ------------------------------------------------------------

// MsgCreateToken defines a message for creating a new token.
type MsgCreateToken struct {
	Creator  string            `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Denom    string            `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Supply   string            `protobuf:"bytes,3,opt,name=supply,proto3" json:"supply,omitempty"` // Use string for input parsing
	Metadata map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty"`
}

// NewMsgCreateToken creates a new MsgCreateToken.
func NewMsgCreateToken(creator string, denom string, supply string, metadata map[string]string) *MsgCreateToken {
	return &MsgCreateToken{
		Creator:  creator,
		Denom:    denom,
		Supply:   supply,
		Metadata: metadata,
	}
}

// Route implements the sdk.Msg interface.
func (MsgCreateToken) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgCreateToken) Type() string { return "CreateToken" }

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return fmt.Errorf("token denom cannot be empty")
	}
	supplyInt, ok := sdk.NewIntFromString(msg.Supply)
	if !ok || supplyInt.IsNil() || supplyInt.IsNegative() {
		return fmt.Errorf("invalid token supply: %s", msg.Supply)
	}
	// Add more validation for metadata if specific fields are required
	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCreateToken) GetSignBytes() sdk.Context {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgCreateToken) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err) // This should have been caught by ValidateBasic
	}
	return []sdk.AccAddress{creator}
}

// RouterKey is the message route for the token module.
const RouterKey = ModuleName

// ModuleCdc is the codec for the token module.
var ModuleCdc = codec.NewProtoCodec(sdk.NewInterfaceRegistry())

// ModuleName defines the module name.
const ModuleName = "token"