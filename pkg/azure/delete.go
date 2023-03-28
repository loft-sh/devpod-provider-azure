package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func deleteVirtualNetWork(ctx context.Context, azureProvider *AzureProvider) error {
	_, exists := checkVirtualNetWork(ctx, azureProvider)
	if !exists {
		return nil
	}

	vnetClient, err := armnetwork.NewVirtualNetworksClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := vnetClient.BeginDelete(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-vnet", nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
func deleteSubnets(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkSubnets(ctx, azureProvider) {
		return nil
	}

	subnetClient, err := armnetwork.NewSubnetsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := subnetClient.BeginDelete(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-vnet",
		azureProvider.Config.MachineID+"-subnet",
		nil,
	)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func deleteNetworkSecurityGroup(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkNetworkSecurityGroup(ctx, azureProvider) {
		return nil
	}

	nsgClient, err := armnetwork.NewSecurityGroupsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := nsgClient.BeginDelete(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-nsg", nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func deletePublicIP(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkPublicIP(ctx, azureProvider) {
		return nil
	}

	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := publicIPAddressClient.BeginDelete(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-public-ip",
		nil,
	)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func deleteNetWorkInterface(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkNetWorkInterface(ctx, azureProvider) {
		return nil
	}

	nicClient, err := armnetwork.NewInterfacesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := nicClient.BeginDelete(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-nic", nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func deleteVirtualMachine(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkVirtualMachine(ctx, azureProvider) {
		return nil
	}

	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := vmClient.BeginDelete(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func deleteDisk(ctx context.Context, azureProvider *AzureProvider) error {
	diskClient, err := armcompute.NewDisksClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}

	pollerResponse, err := diskClient.BeginDelete(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-disk", nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
