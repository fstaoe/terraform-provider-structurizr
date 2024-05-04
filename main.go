// Ensure the documentation is formatted properly
//go:generate terraform fmt -recursive ./examples/

// Generate the documentation for the provider
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name structurizr

package main

import (
	"context"
	"flag"
	"log"

	"github.com/fstaoe/terraform-provider-structurizr/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/fstaoe/structurizr",
		Debug:           debug,
		ProtocolVersion: 6,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
