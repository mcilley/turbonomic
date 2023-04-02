package turbonomic

import (
	"net/url"

	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"
)

// Config struct defines the necessary information needed to configure the
// terraform provider for communication with the Turbonomic API.
type Config struct {
	// Turbonomic URL.  This will be the 'base URL' the REST client uses for issuing
	// API requests to Turbonomic.  This is constructed from the hostname, port, and
	// protocol options passed to the provider.
	Server url.URL
	// Whether or not to verify the server's certificate/hostname.  This flag
	// is passed to the TLS config when initializing the REST client for API
	// communication.
	//
	// See 'pkg/crypto/tls/#Config.InsecureSkipVerify' for more information.
	ClientTLSInsecure bool
	// Set of credentials needed to authenticate against Turbonomic
	ClientCredentials api.ClientCredentials
}

// Creates a client reference for the Turbonomic REST API given the provider
// configuration options.  After creating a client reference, the client
// is then authenticated with the credentials supplied to the provider
// configuration.
func (c *Config) Client() (*api.Client, error) {
	log.Tracef("config.go#Client")

	client := api.NewClient(c.Server, c.ClientTLSInsecure, c.ClientCredentials)

	log.Infof("Rest Client configured")

	return client, nil
}
