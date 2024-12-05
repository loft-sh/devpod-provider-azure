package options

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	AZURE_TAGS            = "AZURE_TAGS"
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
	Tags           map[string]*string
}

type AzureImage struct {
	Offer     string
	Publisher string
	SKU       string
	Version   string
}

func FromEnv(init, withFolder bool) (*Options, error) {
	retOptions := &Options{}

	var err error

	retOptions.ResourceGroup, err = FromEnvOrError(AZURE_RESOURCE_GROUP)
	if err != nil {
		return nil, err
	}

	retOptions.MachineType, err = FromEnvOrError(AZURE_INSTANCE_SIZE)
	if err != nil {
		return nil, err
	}

	image, err := FromEnvOrError(AZURE_IMAGE)
	if err != nil {
		return nil, err
	}
	imageSplit := strings.Split(image, ":")
	if len(imageSplit) < 4 {
		return nil, errors.Errorf("Malformed image name")
	}

	retOptions.DiskImage.Offer = imageSplit[1]
	retOptions.DiskImage.Publisher = imageSplit[0]
	retOptions.DiskImage.SKU = imageSplit[2]
	retOptions.DiskImage.Version = imageSplit[3]

	diskSizeGB, err := FromEnvOrError(AZURE_DISK_SIZE)
	if err != nil {
		return nil, err
	}

	retOptions.DiskSizeGB, err = strconv.Atoi(diskSizeGB)
	if err != nil {
		return nil, err
	}

	retOptions.DiskType, err = FromEnvOrError(AZURE_DISK_TYPE)
	if err != nil {
		return nil, err
	}

	// Optional
	retOptions.CustomData = os.Getenv(AZURE_CUSTOM_DATA)

	retOptions.Zone, err = FromEnvOrError(AZURE_REGION)
	if err != nil {
		return nil, err
	}

	retOptions.SubscriptionID, err = FromEnvOrError(AZURE_SUBSCRIPTION_ID)
	if err != nil {
		return nil, err
	}

	retOptions.Tags, err = parseTags(os.Getenv(AZURE_TAGS))
	if err != nil {
		return nil, err
	}

	// Return eraly if we're just doing init
	if init {
		return retOptions, nil
	}

	retOptions.MachineID, err = FromEnvOrError("MACHINE_ID")
	if err != nil {
		return nil, err
	}
	// prefix with devpod-
	retOptions.MachineID = "devpod-" + retOptions.MachineID

	if withFolder {
		retOptions.MachineFolder, err = FromEnvOrError("MACHINE_FOLDER")
		if err != nil {
			return nil, err
		}
	}

	return retOptions, nil
}

func FromEnvOrError(name string) (string, error) {
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

func parseTags(tagsEnv string) (map[string]*string, error) {
	tags := map[string]*string{}
	if tagsEnv == "" {
		return tags, nil
	}

	tagsRaw := strings.Split(tagsEnv, ",")
	for _, tag := range tagsRaw {
		splitTag := strings.SplitN(tag, "=", 2)
		if len(splitTag) != 2 {
			return tags, fmt.Errorf("Malformed tag, expected format tagName=tagValue: %s", tag)
		}
		tags[splitTag[0]] = to.Ptr[string](splitTag[1])
	}

	return tags, nil
}
