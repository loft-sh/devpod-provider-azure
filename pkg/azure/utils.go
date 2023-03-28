package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/loft-sh/devpod-provider-azure/pkg/options"
	"github.com/loft-sh/devpod/pkg/log"
)

type AzureProvider struct {
	Config           *options.Options
	Cred             *azidentity.DefaultAzureCredential
	Log              log.Logger
	WorkingDirectory string
}

func checkVirtualNetWork(ctx context.Context, azureProvider *AzureProvider) (*armnetwork.VirtualNetwork, bool) {
	vnetClient, err := armnetwork.NewVirtualNetworksClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, false
	}

	resource, err := vnetClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-vnet", nil)
	if err != nil {
		return nil, false
	}

	return &resource.VirtualNetwork, resource.Name != nil
}

func checkSubnets(ctx context.Context, azureProvider *AzureProvider) bool {
	subnetClient, err := armnetwork.NewSubnetsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return false
	}

	resource, err := subnetClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-vnet", azureProvider.Config.MachineID+"-subnet", nil)
	if err != nil {
		return false
	}

	return resource.Name != nil
}

func checkNetworkSecurityGroup(ctx context.Context, azureProvider *AzureProvider) bool {
	nsgClient, err := armnetwork.NewSecurityGroupsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return false
	}

	resource, err := nsgClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-nsg", nil)
	if err != nil {
		return false
	}

	return resource.Name != nil
}

func checkPublicIP(ctx context.Context, azureProvider *AzureProvider) bool {
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return false
	}

	resource, err := publicIPAddressClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-public-ip", nil)
	if err != nil {
		return false
	}

	return resource.Name != nil
}

func checkNetWorkInterface(ctx context.Context, azureProvider *AzureProvider) bool {
	nicClient, err := armnetwork.NewInterfacesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return false
	}

	resource, err := nicClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-nic", nil)
	if err != nil {
		return false
	}

	return resource.Name != nil
}

func checkVirtualMachine(ctx context.Context, azureProvider *AzureProvider) bool {
	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return false
	}

	resource, err := vmClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID, nil)
	if err != nil {
		return false
	}

	return resource.Name != nil
}
