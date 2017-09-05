package credentials

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

type Credentials struct {
	Passwords []Password
	Certificates []Certificate
	SshKeys []SSH
	RsaKeys []RSA
}

func NewPassword(name, value string) Password {
	return Password{
		Metadata: Metadata{
			Base: Base{
				Name: name,
			},
			Type: "password",
		},
		Value: values.Password(value),
	}
}

func NewCertificate(name string, certificate values.Certificate) Certificate {
	return Certificate{
		Metadata: Metadata{
			Base: Base{
				Name: name,
			},
			Type: "certificate",
		},
		Value: certificate,
	}
}

func NewSsh(name string, ssh values.SSH) SSH {
	sshVal := SSH{
		Metadata: Metadata{
			Base: Base{
				Name: name,
			},
			Type: "ssh",
		},
	}
	sshVal.Value.PublicKey = ssh.PublicKey
	sshVal.Value.PrivateKey = ssh.PrivateKey
	return sshVal
}

func NewRsa(name string, rsa values.RSA) RSA {
	return RSA{
		Metadata: Metadata{
			Base: Base{
				Name: name,
			},
			Type: "rsa",
		},
		Value: rsa,
	}
}
