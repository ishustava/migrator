package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

type CredHubClient interface {

	SetPassword(name string, value values.Password, overwrite bool) (credentials.Password, error)

	SetCertificate(name string, value values.Certificate, overwrite bool) (credentials.Certificate, error)

	SetRSA(name string, value values.RSA, overwrite bool) (credentials.RSA, error)

	SetSSH(name string, value values.SSH, overwrite bool) (credentials.SSH, error)
}
