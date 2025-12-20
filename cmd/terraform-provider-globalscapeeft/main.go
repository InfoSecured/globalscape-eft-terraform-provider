package main

import (
	"context"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/provider"
	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/version"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: version.ProviderAddress,
	})
}
