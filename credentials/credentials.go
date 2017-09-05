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

func MakePassword(name, value string) Password {
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

func MakeCertificate(name string, certificate values.Certificate) Certificate {
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

func MakeSsh(name string, ssh values.SSH) SSH {
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

func MakeRsa(name string, rsa values.RSA) RSA {
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
