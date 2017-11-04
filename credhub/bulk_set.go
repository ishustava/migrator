package credhub

import (
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

func BulkSet(credentials *credentials.Credentials, credhubClient credhub.CredHub) error {
	mode := credhub.Overwrite
	for _, pass := range credentials.Passwords {
		credhubClient.SetPassword(pass.Name, pass.Value, mode)
	}

	for _, cert := range credentials.Certificates {
		credhubClient.SetCertificate(cert.Name, cert.Value, mode)
	}

	for _, rsa := range credentials.RsaKeys {
		credhubClient.SetRSA(rsa.Name, rsa.Value, mode)
	}

	for _, ssh := range credentials.SshKeys {
		credhubClient.SetSSH(
			ssh.Name,
			values.SSH{
				PublicKey:  ssh.Value.PublicKey,
				PrivateKey: ssh.Value.PrivateKey,
			},
			mode)
	}

	return nil
}
