package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-azure/pkg/azure"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// TokenCmd holds the cmd flags
type TokenCmd struct{}

// NewTokenCmd defines a command
func NewTokenCmd() *cobra.Command {
	cmd := &TokenCmd{}
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Token an instance",
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

	return tokenCmd
}

// Run runs the command logic
func (cmd *TokenCmd) Run(
	ctx context.Context,
	providerAzure *azure.AzureProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {
	return azure.Token(ctx, providerAzure)
}
