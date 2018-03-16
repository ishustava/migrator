package credhub

import (
	"github.com/ishustava/migrator/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

func BulkSet(credentials *credentials.Credentials, credHubClient CredHubClient, observer BulkSetObserver) error {
	mode := credhub.Overwrite

	observer.BeginBulkSet(
		len(credentials.Passwords),
		len(credentials.Certificates),
		len(credentials.RsaKeys),
		len(credentials.SshKeys),
	)
	for _, pass := range credentials.Passwords {
		if _, err := credHubClient.SetPassword(pass.Name, pass.Value, mode); err != nil {
			observer.FailPasswordSet(pass.Name, err)
		}
	}
	observer.EndPasswordsSet()

	for _, cert := range credentials.Certificates {
		if _, err := credHubClient.SetCertificate(cert.Name, cert.Value, mode); err != nil {
			observer.FailCertificateSet(cert.Name, err)
		}
	}
	observer.EndCertificatesSet()

	for _, rsa := range credentials.RsaKeys {
		if _, err := credHubClient.SetRSA(rsa.Name, rsa.Value, mode); err != nil {
			observer.FailRsaKeySet(rsa.Name, err)
		}
	}
	observer.EndRsaKeysSet()

	for _, ssh := range credentials.SshKeys {
		_, err := credHubClient.SetSSH(
			ssh.Name,
			values.SSH{
				PublicKey:  ssh.Value.PublicKey,
				PrivateKey: ssh.Value.PrivateKey,
			},
			mode)
		if err != nil {
			observer.FailSshKeySet(ssh.Name, err)
		}
	}
	observer.EndSshKeysSet()

	return observer.EndBulkSet()
}