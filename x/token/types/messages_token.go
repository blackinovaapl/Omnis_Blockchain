package types

func NewMsgCreateToken(creator string, name string, symbol string, decimals string, totalSupply string, metadata string) *MsgCreateToken {
	return &MsgCreateToken{
		Creator:     creator,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Metadata:    metadata,
	}
}

func NewMsgUpdateToken(creator string, id uint64, name string, symbol string, decimals string, totalSupply string, metadata string) *MsgUpdateToken {
	return &MsgUpdateToken{
		Id:          id,
		Creator:     creator,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Metadata:    metadata,
	}
}

func NewMsgDeleteToken(creator string, id uint64) *MsgDeleteToken {
	return &MsgDeleteToken{
		Id:      id,
		Creator: creator,
	}
}
