package client

import (
	"context"
	"fmt"
	"strconv"
	
	"github.com/spf13/cobra"

	"omnis/x/token/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// GetTxCmd returns the transaction commands for the token module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateToken())
	cmd.AddCommand(CmdUpdateToken())
	cmd.AddCommand(CmdDeleteToken())

	return cmd
}

// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdShowToken())
	cmd.AddCommand(CmdListTokens())

	return cmd
}

// CmdCreateToken creates a new token.
func CmdCreateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token [name] [symbol] [decimals] [total-supply] [metadata]",
		Short: "Creates a new token with the specified properties",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get the client context from the command
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse the command-line arguments
			name := args[0]
			symbol := args[1]
			decimals, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid decimals: %w", err)
			}
			totalSupply := args[3]
			metadata := args[4]

			// Create the MsgCreateToken transaction message
			msg := types.NewMsgCreateToken(
				clientCtx.GetFromAddress().String(),
				name,
				symbol,
				uint32(decimals),
				totalSupply,
				metadata,
			)

			// Validate the message
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Broadcast the transaction
			return tx.GenerateOrBroadcastTxWithFactory(
				clientCtx,
				clientCtx.TxFactory(),
				msg,
			)
		},
	}

	// Add standard transaction flags (e.g., --from, --gas, --fees, etc.)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateToken updates an existing token.
func CmdUpdateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-token [id] [name] [symbol] [decimals] [total-supply] [metadata]",
		Short: "Updates a token with the new properties",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get the client context from the command
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse the command-line arguments
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token ID: %w", err)
			}
			name := args[1]
			symbol := args[2]
			decimals, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid decimals: %w", err)
			}
			totalSupply := args[4]
			metadata := args[5]

			// Create the MsgUpdateToken transaction message
			msg := types.NewMsgUpdateToken(
				clientCtx.GetFromAddress().String(),
				id,
				name,
				symbol,
				uint32(decimals),
				totalSupply,
				metadata,
			)

			// Validate the message
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Broadcast the transaction
			return tx.GenerateOrBroadcastTxWithFactory(
				clientCtx,
				clientCtx.TxFactory(),
				msg,
			)
		},
	}

	// Add standard transaction flags
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdDeleteToken deletes a token.
func CmdDeleteToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-token [id]",
		Short: "Deletes a token by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get the client context
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse the token ID
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token ID: %w", err)
			}

			// Create the MsgDeleteToken transaction message
			msg := types.NewMsgDeleteToken(
				clientCtx.GetFromAddress().String(),
				id,
			)

			// Validate the message
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Broadcast the transaction
			return tx.GenerateOrBroadcastTxWithFactory(
				clientCtx,
				clientCtx.TxFactory(),
				msg,
			)
		},
	}

	// Add standard transaction flags
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdListTokens lists all tokens.
func CmdListTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tokens",
		Short: "Lists all tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the client context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			
			// Create a query client
			queryClient := types.NewQueryClient(clientCtx)
			
			// Call the ListTokens query endpoint
			res, err := queryClient.Tokens(context.Background(), &types.QueryTokensRequest{})
			if err != nil {
				return err
			}

			// Print the result
			return clientCtx.PrintProto(res)
		},
	}

	// Add standard query flags (e.g., --output, --height)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdShowToken shows a single token.
func CmdShowToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-token [id]",
		Short: "Shows a single token by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get the client context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse the token ID
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token ID: %w", err)
			}
			
			// Create a query client
			queryClient := types.NewQueryClient(clientCtx)
			
			// Call the Token query endpoint with the specified ID
			res, err := queryClient.Token(context.Background(), &types.QueryTokenRequest{Id: id})
			if err != nil {
				return err
			}
			
			// Print the result
			return clientCtx.PrintProto(res)
		},
	}

	// Add standard query flags
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
