package credhub

import (
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

func BulkSet(credentials *credentials.Credentials, credhub CredHubClient) error {
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
