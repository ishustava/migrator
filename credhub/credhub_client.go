package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

type CredHubClient interface {
	SetPassword(name string, value values.Password, overwrite credhub.Mode) (credentials.Password, error)
	SetCertificate(name string, value values.Certificate, overwrite credhub.Mode) (credentials.Certificate, error)
	SetRSA(name string, value values.RSA, overwrite credhub.Mode) (credentials.RSA, error)
	SetSSH(name string, value values.SSH, overwrite credhub.Mode) (credentials.SSH, error)
}
