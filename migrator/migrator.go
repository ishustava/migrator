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
	}

	return creds, nil
}

func tryUnmarshalCertificate(value interface{}) (values.Certificate, error) {
	certYaml, _ := yaml.Marshal(value)

	certificate := values.Certificate{}
	err := yaml.Unmarshal(certYaml, &certificate)

	return certificate, err
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
