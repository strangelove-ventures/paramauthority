package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkupgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/strangelove-ventures/paramauthority/x/upgrade/types"
)

const (
	FlagUpgradeHeight = "upgrade-height"
	FlagUpgradeInfo   = "upgrade-info"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   sdkupgradetypes.ModuleName,
		Short: "Upgrade transaction subcommands",
	}

	cmd.AddCommand(NewCmdSubmitUpgrade())
	cmd.AddCommand(NewCmdSubmitCancelUpgrade())

	return cmd
}

// NewCmdSubmitUpgrade implements a command handler for submitting a software upgrade transaction.
func NewCmdSubmitUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software-upgrade [name] (--upgrade-height [height]) (--upgrade-info [info]) [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a software upgrade",
		Long: "Submit a software upgrade along with an initial deposit.\n" +
			"Please specify a unique name and height for the upgrade to take effect.\n" +
			"You may include info to reference a binary download link, in a format compatible with: https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			plan, err := parseArgsToPlan(cmd, name)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg := &types.MsgSoftwareUpgrade{
				Authority: from.String(),
				Plan:      plan,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Int64(FlagUpgradeHeight, 0, "The height at which the upgrade must happen")
	cmd.Flags().String(FlagUpgradeInfo, "", "Optional info for the planned upgrade such as commit hash, etc.")

	return cmd
}

// NewCmdSubmitCancelUpgrade implements a command handler for submitting a software upgrade cancel transaction.
func NewCmdSubmitCancelUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-software-upgrade [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Cancel the current software upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			msg := &types.MsgCancelUpgrade{
				Authority: from.String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func parseArgsToPlan(cmd *cobra.Command, name string) (sdkupgradetypes.Plan, error) {
	height, err := cmd.Flags().GetInt64(FlagUpgradeHeight)
	if err != nil {
		return sdkupgradetypes.Plan{}, err
	}

	info, err := cmd.Flags().GetString(FlagUpgradeInfo)
	if err != nil {
		return sdkupgradetypes.Plan{}, err
	}

	return sdkupgradetypes.Plan{Name: name, Height: height, Info: info}, nil
}
