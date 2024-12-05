package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-azure/pkg/azure"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// StopRemoteCmd holds the cmd flags
type StopRemoteCmd struct{}

// NewStopRemoteCmd defines a command
func NewStopRemoteCmd() *cobra.Command {
	cmd := &StopRemoteCmd{}
	stopRemoteCmd := &cobra.Command{
		Use:   "stop-remote",
		Short: "StopRemote an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			azureProvider, err := azure.NewProvider(false, log.Default)
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

	return stopRemoteCmd
}

// Run runs the command logic
func (cmd *StopRemoteCmd) Run(
	ctx context.Context,
	providerAzure *azure.AzureProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {
	return azure.StopRemote(ctx, providerAzure)
}
