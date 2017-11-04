package credhub

import (
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

func BulkSet(credentials *credentials.Credentials, credHubClient CredHubClient) error {
	mode := credhub.Overwrite
	for _, pass := range credentials.Passwords {
		credHubClient.SetPassword(pass.Name, pass.Value, mode)
	}

	for _, cert := range credentials.Certificates {
		credHubClient.SetCertificate(cert.Name, cert.Value, mode)
	}

	for _, rsa := range credentials.RsaKeys {
		credHubClient.SetRSA(rsa.Name, rsa.Value, mode)
	}

	for _, ssh := range credentials.SshKeys {
		credHubClient.SetSSH(
			ssh.Name,
			values.SSH{
				PublicKey: ssh.Value.PublicKey,
				PrivateKey: ssh.Value.PrivateKey,
			},
			mode)
	}

	return nil
}
