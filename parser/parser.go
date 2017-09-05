package parser

import (
	"errors"
	"gopkg.in/yaml.v2"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/ishustava/migrator/credentials"
)

func FindCredentials(varsStore map[string]interface{}) (*Credentials, error) {
	creds := new(Credentials)

	for name, value := range varsStore {
		passwordValue, ok := value.(string)
		if ok {
			creds.Passwords = append(creds.Passwords, NewPassword(name, passwordValue))
			continue
		}

		certificateValue, err := tryUnmarshalCertificate(value)
		if err == nil {
			creds.Certificates = append(creds.Certificates, NewCertificate(name, certificateValue))
			continue
		}

		sshValue, err := tryUnmarshalSsh(value)
		if err == nil {
			creds.SshKeys = append(creds.SshKeys, NewSsh(name, sshValue))
			continue
		}

		rsaValue, err := tryUnmarshalRsa(value)
		if err == nil {
			creds.RsaKeys = append(creds.RsaKeys, NewRsa(name, rsaValue))
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

