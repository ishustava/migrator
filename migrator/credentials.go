package migrator

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

type Credentials struct {
	Passwords []credentials.Password
	Certificates []credentials.Certificate
	SshKeys []credentials.SSH
	RsaKeys []credentials.RSA
}
