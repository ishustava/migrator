package migrator

import (
	"io/ioutil"
	"errors"
	"gopkg.in/yaml.v2"
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

func parseVarsStoreFile(path string) (*Credentials, error) {
	creds := new(Credentials)

	varsStoreYaml, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error occurred when reading vars store file: " + err.Error())
	}

	varsStore := make(map[string]interface{})

	err = yaml.Unmarshal(varsStoreYaml, &varsStore)
	if err != nil {
		return nil, errors.New("Could not unmarshal vars store file: " + err.Error())
	}

	for name, value := range varsStore {
		passwordValue, ok := value.(string)
		if ok {
			creds.Passwords = append(creds.Passwords, makePassword(name, passwordValue))
			continue
		}

		certificateValue, err := tryUnmarshalCertificate(value)
		if err == nil {
			creds.Certificates = append(creds.Certificates, makeCertificate(name, certificateValue))
			continue
		}

		sshValue, err := tryUnmarshalSsh(value)
		if err == nil {
			creds.SshKeys = append(creds.SshKeys, makeSsh(name, sshValue))
			continue
		}

		rsaValue, err := tryUnmarshalRsa(value)
		if err == nil {
			creds.RsaKeys = append(creds.RsaKeys, makeRsa(name, rsaValue))
			continue
		}
	}

	return creds, nil
}

func tryUnmarshalCertificate(value interface{}) (values.Certificate, error) {
	certYaml, _ := yaml.Marshal(value)

	certificate := values.Certificate{}
	err := yaml.UnmarshalStrict(certYaml, &certificate)

	return certificate, err
}

func tryUnmarshalSsh(value interface{}) (values.SSH, error) {
	ssh := values.SSH{}

	sshMap := value.(map[interface{}]interface{})
	_, ok := sshMap["public_key_fingerprint"]
	if !ok {
		return ssh, errors.New("Key not found in map: public_key_fingerprint")
	}
	delete(sshMap, "public_key_fingerprint")

	sshYaml, _ := yaml.Marshal(value)

	err := yaml.UnmarshalStrict(sshYaml, &ssh)

	return ssh, err
}

func tryUnmarshalRsa(value interface{}) (values.RSA, error) {
	rsaYaml, _ := yaml.Marshal(value)

	rsa := values.RSA{}
	err := yaml.UnmarshalStrict(rsaYaml, &rsa)

	return rsa, err
}

func makePassword(name, value string) Password {
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

func makeCertificate(name string, certificate values.Certificate) Certificate {
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

func makeSsh(name string, ssh values.SSH) SSH {
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

func makeRsa(name string, rsa values.RSA) RSA {
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
