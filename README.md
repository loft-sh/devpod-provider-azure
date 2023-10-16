# AZURE Provider for DevPod

[![Join us on Slack!](docs/static/media/slack.svg)](https://slack.loft.sh/) [![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/loft-sh/devpod-provider-azure)

## Getting started

The provider is available for auto-installation using 

```sh
devpod provider add azure
devpod provider use azure
```

Follow the on-screen instructions to complete the setup.

Needed variables will be:

- AZURE_RESOURCE_GROUP
- AZURE_REGION

### Creating your first devpod env with azure

After the initial setup, just use:

```sh
devpod up .
```

You'll need to wait for the machine and environment setup.

Be aware that authentication is obtained using Azure's Default Credential authenticator, this uses
the CLI tool, the ENV or Certificates, take a look
[here](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli)
for more info on how to setup either one of those auth methods.


### Customize the VM Instance

This provides has the seguent options

|    NAME           | REQUIRED |          DESCRIPTION                  |         DEFAULT         |
|-------------------|----------|---------------------------------------|-------------------------|
| AZURE_DISK_SIZE           | false    | The disk size to use.          | 40                                      |
| AZURE_IMAGE               | false    | The disk image to use.         | Canonical:UbuntuServer:18.04-LTS:latest |
| AZURE_INSTANCE_SIZE       | false    | The machine type to use.       | Standard_D11_v2                         |
| AZURE_REGION              | true     | The azure region to use        |                                         |
| AZURE_RESOURCE_GROUP      | true     | The azure resource group name  |                                         |
| AZURE_SUBSCRIPTION_ID     | true     | The azure subscription id      |                                         |

Options can either be set in `env` or using for example:

```sh
devpod provider set-options -o AZURE_IMAGE=Vendor:Item:Version:Tag
```
