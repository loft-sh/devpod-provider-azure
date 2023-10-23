package azure

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/loft-sh/devpod/pkg/ssh"
)

func createVirtualNetwork(ctx context.Context, azureProvider *AzureProvider) (*armnetwork.VirtualNetwork, error) {
	vnet, exists := checkVirtualNetWork(ctx, azureProvider)
	if exists {
		return vnet, nil
	}

	vnetClient, err := armnetwork.NewVirtualNetworksClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.VirtualNetwork{
		Location: to.Ptr(azureProvider.Config.Zone),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16"), // example 10.1.0.0/16
				},
			},
		},
	}

	pollerResponse, err := vnetClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-vnet",
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualNetwork, nil
}

func createSubnets(ctx context.Context, azureProvider *AzureProvider) (*armnetwork.Subnet, error) {
	subnetClient, err := armnetwork.NewSubnetsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.Ptr("10.1.10.0/24"),
		},
	}

	pollerResponse, err := subnetClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-vnet",
		azureProvider.Config.MachineID+"-subnet",
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Subnet, nil
}

func createNetworkSecurityGroup(
	ctx context.Context,
	azureProvider *AzureProvider,
) (*armnetwork.SecurityGroup, error) {
	nsgClient, err := armnetwork.NewSecurityGroupsClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.SecurityGroup{
		Location: to.Ptr(azureProvider.Config.Zone),
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			SecurityRules: []*armnetwork.SecurityRule{
				// Windows connection to virtual machine needs to open port 3389,RDP
				// inbound
				{
					Name: to.Ptr("devpod_inbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description: to.Ptr(
							"devpod network security group inbound port 22",
						),
						Direction: to.Ptr(armnetwork.SecurityRuleDirectionInbound),
					},
				},
				// outbound
				{
					Name: to.Ptr("devpod_outbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description: to.Ptr(
							"devpod network security group outbound port 22",
						),
						Direction: to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
					},
				},
			},
		},
	}

	pollerResponse, err := nsgClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-nsg",
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SecurityGroup, nil
}

func createPublicIP(
	ctx context.Context,
	azureProvider *AzureProvider,
) (*armnetwork.PublicIPAddress, error) {
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.PublicIPAddress{
		Location: to.Ptr(azureProvider.Config.Zone),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: to.Ptr(
				armnetwork.IPAllocationMethodStatic,
			), // Static or Dynamic
		},
	}

	pollerResponse, err := publicIPAddressClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-public-ip",
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.PublicIPAddress, err
}

func createNetWorkInterface(
	ctx context.Context,
	azureProvider *AzureProvider,

	subnetID string,
	publicIPID string,
	networkSecurityGroupID string,
) (*armnetwork.Interface, error) {
	nicClient, err := armnetwork.NewInterfacesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.Interface{
		Location: to.Ptr(azureProvider.Config.Zone),
		Properties: &armnetwork.InterfacePropertiesFormat{
			// NetworkSecurityGroup:
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr("ipConfig"),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						Subnet: &armnetwork.Subnet{
							ID: to.Ptr(subnetID),
						},
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: to.Ptr(publicIPID),
						},
					},
				},
			},
			NetworkSecurityGroup: &armnetwork.SecurityGroup{
				ID: to.Ptr(networkSecurityGroupID),
			},
		},
	}

	pollerResponse, err := nicClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID+"-nic",
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Interface, err
}

func createVirtualMachine(
	ctx context.Context,
	azureProvider *AzureProvider,
	networkInterfaceID string,
) (*armcompute.VirtualMachine, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return nil, err
	}

	publicKeyBase, err := ssh.GetPublicKeyBase(azureProvider.Config.MachineFolder)
	if err != nil {
		return nil, err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return nil, err
	}

	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(azureProvider.Config.Zone),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					// search image reference
					// az vm image list --output table
					Offer:     to.Ptr(azureProvider.Config.DiskImage.Offer),
					Publisher: to.Ptr(azureProvider.Config.DiskImage.Publisher),
					SKU:       to.Ptr(azureProvider.Config.DiskImage.SKU),
					Version:   to.Ptr(azureProvider.Config.DiskImage.Version),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(azureProvider.Config.MachineID + "-disk"),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(
							armcompute.StorageAccountTypes(azureProvider.Config.DiskType),
						), // OSDisk type Standard/Premium HDD/SSD
					},
					DiskSizeGB: to.Ptr[int32](int32(azureProvider.Config.DiskSizeGB)),
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(
					armcompute.VirtualMachineSizeTypes(azureProvider.Config.MachineType),
				), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: &armcompute.OSProfile{ //
				ComputerName:  to.Ptr(azureProvider.Config.MachineID),
				AdminUsername: to.Ptr("devpod"),
				CustomData:    to.Ptr(azureProvider.Config.CustomData),
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{
							{
								Path: to.Ptr(
									fmt.Sprintf("/home/%s/.ssh/authorized_keys", "devpod"),
								),
								KeyData: to.Ptr(string(publicKey)),
							},
						},
					},
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(networkInterfaceID),
					},
				},
			},
		},
	}

	pollerResponse, err := vmClient.BeginCreateOrUpdate(
		ctx,
		azureProvider.Config.ResourceGroup,
		azureProvider.Config.MachineID,
		parameters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualMachine, nil
}
