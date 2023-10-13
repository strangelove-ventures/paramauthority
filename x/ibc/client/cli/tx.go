package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v4/modules/core/exported"
	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/paramauthority/x/ibc/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "IBC authority transaction subcommands",
	}

	cmd.AddCommand(NewCmdSubmitUpdateClientProposal())
	cmd.AddCommand(NewCmdSubmitUpgradeProposal())

	return cmd
}

// NewCmdSubmitUpdateClientProposal implements a command handler for submitting an update IBC client proposal transaction.
func NewCmdSubmitUpdateClientProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client [subject-client-id] [substitute-client-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit an update IBC client proposal",
		Long: "Submit an update IBC client\n" +
			"Please specify a subject client identifier you want to update..\n" +
			"Please specify the substitute client the subject client will be updated to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subjectClientID, substituteClientID := args[0], args[1]

			from := clientCtx.GetFromAddress()

			msg := &types.MsgClientUpdate{
				Authority:          from.String(),
				SubjectClientId:    subjectClientID,
				SubstituteClientId: substituteClientID,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewCmdSubmitUpgradeProposal implements a command handler for submitting an upgrade IBC client proposal transaction.
func NewCmdSubmitUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-upgrade [name] [height] [path/to/upgraded_client_state.json] [flags]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit an IBC upgrade proposal",
		Long: "Submit an IBC client breaking upgrade proposal along with an initial deposit.\n" +
			"The client state specified is the upgraded client state representing the upgraded chain\n" +
			`Example Upgraded Client State JSON: 
{
	"@type":"/ibc.lightclients.tendermint.v1.ClientState",
 	"chain_id":"testchain1",
	"unbonding_period":"1814400s",
	"latest_height":{"revision_number":"0","revision_height":"2"},
	"proof_specs":[{"leaf_spec":{"hash":"SHA256","prehash_key":"NO_HASH","prehash_value":"SHA256","length":"VAR_PROTO","prefix":"AA=="},"inner_spec":{"child_order":[0,1],"child_size":33,"min_prefix_length":4,"max_prefix_length":12,"empty_child":null,"hash":"SHA256"},"max_depth":0,"min_depth":0},{"leaf_spec":{"hash":"SHA256","prehash_key":"NO_HASH","prehash_value":"SHA256","length":"VAR_PROTO","prefix":"AA=="},"inner_spec":{"child_order":[0,1],"child_size":32,"min_prefix_length":1,"max_prefix_length":1,"empty_child":null,"hash":"SHA256"},"max_depth":0,"min_depth":0}],
	"upgrade_path":["upgrade","upgradedIBCState"],
}
			`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			name := args[0]

			height, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			plan := upgradetypes.Plan{
				Name:   name,
				Height: height,
			}

			// attempt to unmarshal client state argument
			var clientState exported.ClientState
			clientContentOrFileName := args[2]
			if err := cdc.UnmarshalInterfaceJSON([]byte(clientContentOrFileName), &clientState); err != nil {

				// check for file path if JSON input is not provided
				contents, err := os.ReadFile(clientContentOrFileName)
				if err != nil {
					return fmt.Errorf("neither JSON input nor path to .json file for client state were provided: %w", err)
				}

				if err := cdc.UnmarshalInterfaceJSON(contents, &clientState); err != nil {
					return fmt.Errorf("error unmarshalling client state file: %w", err)
				}
			}

			from := clientCtx.GetFromAddress()

			any, err := clienttypes.PackClientState(clientState)
			if err != nil {
				return err
			}

			msg := &types.MsgUpgrade{
				Authority:           from.String(),
				Plan:                plan,
				UpgradedClientState: any,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
