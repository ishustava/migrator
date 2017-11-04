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
	VarsStore           string   `short:"v" long:"vars-store" required:"yes" description:"Path to vars store file."`
	CredHubURL          string   `short:"u" long:"credhub-url" description:"URL of the CredHub server, e.g. https://example.com:8844" env:"CREDHUB_URL"`
	CredHubClient       string   `short:"c" long:"credhub-client" description:"UAA client for the CredHub Server" env:"CREDHUB_CLIENT"`
	CredHubClientSecret string   `short:"s" long:"credhub-client-secret" description:"UAA client secret for the CredHub Server" env:"CREDHUB_CLIENT_SECRET"`
	DirectorName        string   `short:"e" long:"director-name" description:"Name of the BOSH director"`
	DeploymentName      string   `short:"d" long:"deployment-name" description:"Name of the BOSH deployment with which vars store is used"`
	CaCerts             []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections. Multiple flags may be provided." env:"CREDHUB_CA_CERT"`
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

	credentials, err := parser.ParseCredentials(varsStore)
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
