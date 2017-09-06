package migrator

import (
	"github.com/ishustava/migrator/credentials"
	credhub_client "github.com/ishustava/migrator/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

func Migrate(credentials credentials.Credentials, credhub credhub_client.CredHubClient) error {
	for _, pass := range credentials.Passwords {
		credhub.SetPassword(pass.Name, pass.Value, true)
	}

	for _, cert := range credentials.Certificates {
		credhub.SetCertificate(cert.Name, cert.Value, true)
	}

	for _, rsa := range credentials.RsaKeys {
		credhub.SetRSA(rsa.Name, rsa.Value, true)
	}

	for _, ssh := range credentials.SshKeys {
		credhub.SetSSH(
			ssh.Name,
			values.SSH{
				PublicKey: ssh.Value.PublicKey,
				PrivateKey: ssh.Value.PrivateKey,
			},
			true)
	}

	return nil
}
