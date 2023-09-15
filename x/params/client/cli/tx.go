package cli

import (
	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/paramauthority/x/params/types/proposal"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkparamstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	sdkproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   sdkparamstypes.ModuleName,
		Short: "Params transaction subcommands",
	}

	cmd.AddCommand(NewCmdUpdateParams())

	return cmd
}

// NewCmdSubmitUpgrade implements a command handler for submitting a software upgrade transaction.
func NewCmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [subspace] [key] [value]",
		Args:  cobra.ExactArgs(3),
		Short: "Update subspace params",
		Long:  "Update the value for a key in a params subspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			subspace, key, value := args[0], args[1], args[2]

			from := clientCtx.GetFromAddress()

			msg := &proposal.MsgUpdateParams{
				Authority: from.String(),
				ChangeProposal: &sdkproposal.ParameterChangeProposal{
					Changes: []sdkproposal.ParamChange{
						{
							Subspace: subspace,
							Key:      key,
							Value:    value,
						},
					},
				},
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
