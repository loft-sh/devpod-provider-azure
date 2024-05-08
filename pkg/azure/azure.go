package azure

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	"github.com/loft-sh/devpod-provider-azure/pkg/options"
)

func NewProvider(logs log.Logger) (*AzureProvider, error) {
	config, err := options.FromEnv(false)
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.Errorf("Authentication failure: %+v", err)
	}

	// create provider
	provider := &AzureProvider{
		Config: config,
		Cred:   cred,
		Log:    logs,
	}

	return provider, nil
}

func Create(ctx context.Context, azureProvider *AzureProvider) error {
	_, err := createVirtualNetwork(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot create virtual network:%+v", err)
	}

	subnet, err := createSubnets(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot create subnet:%+v", err)
	}

	publicIP, err := createPublicIP(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot create public IP address:%+v", err)
	}

	// network security group
	nsg, err := createNetworkSecurityGroup(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot create network security group:%+v", err)
	}

	netWorkInterface, err := createNetWorkInterface(
		ctx,
		azureProvider,
		*subnet.ID,
		*publicIP.ID,
		*nsg.ID,
	)
	if err != nil {
		return errors.Errorf("cannot create network interface:%+v", err)
	}

	networkInterfaceID := netWorkInterface.ID
	_, err = createVirtualMachine(
		ctx,
		azureProvider,
		*networkInterfaceID,
	)
	if err != nil {
		return errors.Errorf("cannot create virual machine:%+v", err)
	}

	return nil
}

func Delete(ctx context.Context, azureProvider *AzureProvider) error {
	err := deleteVirtualMachine(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete virtual machine:%+v", err)
	}

	err = deleteDisk(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete disk:%+v", err)
	}

	err = deleteNetWorkInterface(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete network interface:%+v", err)
	}

	err = deleteNetworkSecurityGroup(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete network security group:%+v", err)
	}

	err = deletePublicIP(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete public IP address:%+v", err)
	}

	err = deleteSubnets(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete subnet:%+v", err)
	}

	err = deleteVirtualNetWork(ctx, azureProvider)
	if err != nil {
		return errors.Errorf("cannot delete virtual network:%+v", err)
	}

	return nil
}

func Status(ctx context.Context, azureProvider *AzureProvider) (client.Status, error) {
	if !checkVirtualMachine(ctx, azureProvider) {
		return client.StatusNotFound, nil
	}

	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return client.StatusNotFound, nil
	}

	resource, err := vmClient.InstanceView(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID, nil)
	if err != nil {
		return client.StatusNotFound, nil
	}

	status := resource.Statuses[1].DisplayStatus

	switch {
	case *status == "VM running":
		return client.StatusRunning, nil
	case *status == "VM deallocated":
		return client.StatusStopped, nil
	case *status == "VM stopped":
		return client.StatusStopped, nil
	default:
		return client.StatusBusy, nil
	}
}

func Stop(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkVirtualMachine(ctx, azureProvider) {
		return nil
	}

	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}
	pollerResponse, err := vmClient.BeginDeallocate(
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

// we fallback to use direct http request for this as we'll be on the remote machine.
func StopRemote(ctx context.Context, azureProvider *AzureProvider) error {
	client := &http.Client{}

	data := strings.NewReader(``)

	token, err := options.FromEnvOrError("AZURE_PROVIDER_TOKEN")
	if err != nil {
		return err
	}

	url := "https://management.azure.com/subscriptions/" +
		azureProvider.Config.SubscriptionID + "/resourceGroups/" +
		azureProvider.Config.ResourceGroup + "/providers/Microsoft.Compute/virtualMachines/" +
		azureProvider.Config.MachineID + "/deallocate?api-version=2024-03-01"

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func Token(ctx context.Context, azureProvider *AzureProvider) error {
	out, err := exec.Command("az", []string{"account", "get-access-token", "--query", "accessToken"}...).Output()
	if err != nil {
		return err
	}

	token := strings.ReplaceAll(string(out), "\"", "")
	token = strings.Trim(token, "\n")

	fmt.Println(token)
	return nil
}

func Start(ctx context.Context, azureProvider *AzureProvider) error {
	if !checkVirtualMachine(ctx, azureProvider) {
		return nil
	}

	vmClient, err := armcompute.NewVirtualMachinesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return err
	}
	pollerResponse, err := vmClient.BeginStart(
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

func GetInstanceIP(ctx context.Context, azureProvider *AzureProvider) (string, error) {
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(azureProvider.Config.SubscriptionID, azureProvider.Cred, nil)
	if err != nil {
		return "", err
	}

	resource, err := publicIPAddressClient.Get(ctx, azureProvider.Config.ResourceGroup, azureProvider.Config.MachineID+"-public-ip", nil)
	if err != nil {
		return "", err
	}

	return *resource.PublicIPAddress.Properties.IPAddress, nil
}
