package options

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	AZURE_REGION          = "AZURE_REGION"
	AZURE_INSTANCE_SIZE   = "AZURE_INSTANCE_SIZE"
	AZURE_IMAGE           = "AZURE_IMAGE"
	AZURE_RESOURCE_GROUP  = "AZURE_RESOURCE_GROUP"
	AZURE_DISK_TYPE       = "AZURE_DISK_TYPE"
	AZURE_DISK_SIZE       = "AZURE_DISK_SIZE"
	AZURE_CUSTOM_DATA     = "AZURE_CUSTOM_DATA"
	AZURE_SUBSCRIPTION_ID = "AZURE_SUBSCRIPTION_ID"
)

type Options struct {
	DiskImage      AzureImage
	DiskSizeGB     int
	DiskType       string
	CustomData     string
	MachineFolder  string
	MachineID      string
	MachineType    string
	ResourceGroup  string
	SubscriptionID string
	Zone           string
}

type AzureImage struct {
	Offer     string
	Publisher string
	SKU       string
	Version   string
}

func FromEnv(init bool) (*Options, error) {
	retOptions := &Options{}

	var err error

	retOptions.ResourceGroup, err = fromEnvOrError(AZURE_RESOURCE_GROUP)
	if err != nil {
		return nil, err
	}

	retOptions.MachineType, err = fromEnvOrError(AZURE_INSTANCE_SIZE)
	if err != nil {
		return nil, err
	}

	image, err := fromEnvOrError(AZURE_IMAGE)
	if err != nil {
		return nil, err
	}
	imageSplit := strings.Split(image, ":")
	if len(imageSplit) < 4 {
		return nil, errors.Errorf("Malformet image name")
	}

	retOptions.DiskImage.Offer = imageSplit[1]
	retOptions.DiskImage.Publisher = imageSplit[0]
	retOptions.DiskImage.SKU = imageSplit[2]
	retOptions.DiskImage.Version = imageSplit[3]

	diskSizeGB, err := fromEnvOrError(AZURE_DISK_SIZE)
	if err != nil {
		return nil, err
	}

	retOptions.DiskSizeGB, err = strconv.Atoi(diskSizeGB)
	if err != nil {
		return nil, err
	}

	retOptions.DiskType, err = fromEnvOrError(AZURE_DISK_TYPE)
	if err != nil {
		return nil, err
	}

	retOptions.CustomData, err = fromEnvOrError(AZURE_CUSTOM_DATA)
	if err != nil {
		return nil, err
	}

	retOptions.Zone, err = fromEnvOrError(AZURE_REGION)
	if err != nil {
		return nil, err
	}

	retOptions.SubscriptionID, err = fromEnvOrError(AZURE_SUBSCRIPTION_ID)
	if err != nil {
		return nil, err
	}

	// Return eraly if we're just doing init
	if init {
		return retOptions, nil
	}

	retOptions.MachineID, err = fromEnvOrError("MACHINE_ID")
	if err != nil {
		return nil, err
	}
	// prefix with devpod-
	retOptions.MachineID = "devpod-" + retOptions.MachineID

	retOptions.MachineFolder, err = fromEnvOrError("MACHINE_FOLDER")
	if err != nil {
		return nil, err
	}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf(
			"couldn't find option %s in environment, please make sure %s is defined",
			name,
			name,
		)
	}

	return val, nil
}
