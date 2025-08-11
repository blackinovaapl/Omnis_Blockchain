package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		TokenList: []Token{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	tokenIdMap := make(map[uint64]bool)
	tokenCount := gs.GetTokenCount()
	for _, elem := range gs.TokenList {
		if _, ok := tokenIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for token")
		}
		if elem.Id >= tokenCount {
			return fmt.Errorf("token id should be lower or equal than the last id")
		}
		tokenIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
