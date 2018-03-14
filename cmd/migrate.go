package cmd

import (
	"fmt"
	credhubClient "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/commands"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/ishustava/migrator/parser"
	"github.com/ishustava/migrator/credhub"
)

type MigrateCommand struct {
	VarsStore           string   `short:"v" long:"vars-store" required:"yes" description:"Path to vars store file." required:"true"`
	CredHubURL          string   `short:"u" long:"credhub-server" description:"URL of the CredHub server, e.g. https://example.com:8844" env:"CREDHUB_SERVER" required:"true"`
	CredHubClient       string   `short:"c" long:"credhub-client" description:"UAA client for the CredHub Server" env:"CREDHUB_CLIENT" required:"true"`
	CredHubClientSecret string   `short:"s" long:"credhub-secret" description:"UAA client secret for the CredHub Server" env:"CREDHUB_SECRET" required:"true"`
	DirectorName        string   `short:"e" long:"director-name" description:"Name of the BOSH director" required:"true"`
	DeploymentName      string   `short:"d" long:"deployment-name" description:"Name of the BOSH deployment with which vars store is used" env:"BOSH_DEPLOYMENT" required:"true"`
	CaCerts             []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections. Multiple flags may be provided." env:"CREDHUB_CA_CERT" required:"true"`
}

func (cmd MigrateCommand) Execute([]string) error {
	varsStoreFileContents, err := ioutil.ReadFile(cmd.VarsStore)
	if err != nil {
		return err
	}

	var varsStore map[string]interface{}

	err = yaml.Unmarshal(varsStoreFileContents, &varsStore)
	if err != nil {
		return err
	}

	varsStore = parser.AddBoshNamespacing(varsStore, cmd.DirectorName, cmd.DeploymentName)
	credentials, err := parser.ParseCredentials(varsStore)
	if err != nil {
		return err
	}

	credentials.Certificates, err = parser.FindAndSetSigningCA(credentials.Certificates)
	if err != nil {
		return err
	}

	caCerts, err := commands.ReadOrGetCaCerts(cmd.CaCerts)
	if err != nil {
		return err
	}

	ch, err := credhubClient.New(
		cmd.CredHubURL,
		credhubClient.CaCerts(caCerts...),
		credhubClient.Auth(
			auth.UaaClientCredentials(cmd.CredHubClient, cmd.CredHubClientSecret)),
	)
	if err != nil {
		return err
	}

	err = credhub.BulkSet(credentials, ch)
	if err == nil {
		fmt.Println("Successfully migrated all credentials")
	}

	return err
}
