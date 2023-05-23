package api

type Client interface {
	HypervDvdClient
	HypervVhdClient
	HypervVmClient
	HypervVmDvdDriveClient
	HypervVmFirmwareClient
	HypervVmHardDiskDriveClient
	HypervVmIntegrationServiceClient
	HypervVmNetworkAdapterClient
	HypervVmProcessorClient
	HypervVmStatusClient
	HypervVmSwitchClient
}

type Provider struct {
	Client Client
}
