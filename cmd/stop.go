package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-azure/pkg/azure"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct{}

// NewStopCmd defines a command
func NewStopCmd() *cobra.Command {
	cmd := &StopCmd{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			azureProvider, err := azure.NewProvider(log.Default)
			if err != nil {
				return err
			}

			return cmd.Run(
				context.Background(),
				azureProvider,
				provider.FromEnvironment(),
				log.Default,
			)
		},
	}

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(
	ctx context.Context,
	providerAzure *azure.AzureProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {
	return azure.Stop(ctx, providerAzure)
}
